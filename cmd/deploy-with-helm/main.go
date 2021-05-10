package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"sigs.k8s.io/yaml"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func main() {
	// optional flags (can be empty or not given)
	namespaceFlag := flag.String("namespace", "", "namespace")
	projectFlag := flag.String("project", "", "project")
	repositoryFlag := flag.String("repository", "", "repository")
	componentFlag := flag.String("component", "", "component")
	gitRefSpecFlag := flag.String("git-ref-spec", "", "Git ref spec")
	gitCommitSHAFlag := flag.String("git-commit-sha", "", "Git commit SHA")
	releaseNameFlag := flag.String("release-name", "", "release-name")

	// required flags (but not needed as task input)
	nexusURLFlag := flag.String("nexus-url", os.Getenv("NEXUS_URL"), "Nexus URL")
	nexusUsernameFlag := flag.String("nexus-username", os.Getenv("NEXUS_USERNAME"), "Nexus username")
	nexusPasswordFlag := flag.String("nexus-password", os.Getenv("NEXUS_PASSWORD"), "Nexus password")
	chartDir := flag.String("chart-dir", "", "Chart dir")
	environment := flag.String("environment", "", "environment")
	target := flag.String("target", "", "target")
	flag.Parse()

	var namespace string
	if len(*namespaceFlag) > 0 {
		namespace = *namespaceFlag
	} else {
		kubernetesNamespace, err := getTrimmedFileContent(namespaceFile)
		if err != nil {
			panic(err.Error())
		}
		namespace = kubernetesNamespace
	}
	fmt.Printf("namespace=%s\n", namespace)

	var project string
	if len(*projectFlag) > 0 {
		project = *projectFlag
	} else {
		project = strings.TrimSuffix(namespace, "-cd")
	}
	fmt.Printf("project=%s\n", project)

	var repository string
	if len(*repositoryFlag) > 0 {
		repository = *repositoryFlag
	} else {
		r, err := getGitRepository()
		check(err)
		repository = r
	}
	fmt.Printf("repository=%s\n", repository)

	var component string
	if len(*componentFlag) > 0 {
		component = *componentFlag
	} else {
		component = strings.TrimPrefix(repository, fmt.Sprintf("%s-", project))
	}
	fmt.Printf("component=%s\n", component)

	var releaseName string
	if len(*releaseNameFlag) > 0 {
		releaseName = *releaseNameFlag
	} else {
		releaseName = component
	}
	fmt.Printf("releaseName=%s\n", releaseName)

	var gitRefSpec string
	if len(*gitRefSpecFlag) > 0 {
		gitRefSpec = *gitRefSpecFlag
	} else {
		grs, err := getGitFullRef()
		check(err)
		gitRefSpec = grs
	}
	fmt.Printf("gitRefSpec=%s\n", gitRefSpec)

	var gitCommitSHA string
	if len(*gitCommitSHAFlag) > 0 {
		gitCommitSHA = *gitCommitSHAFlag
	} else {
		gcs, err := getGitCommitSHA(gitRefSpec)
		check(err)
		gitCommitSHA = gcs
	}
	fmt.Printf("gitCommitSHA=%s\n", gitCommitSHA)

	// read ods.yml
	odsConfig, err := getConfig("ods.yml")
	if err != nil {
		panic(fmt.Sprintf("err during ods config reading %s", err))
	}
	targetConfig, err := getTarget(odsConfig, *environment, *target)
	if err != nil {
		panic(fmt.Sprintf("err during namespace extraction %s", err))
	}

	releaseNamespace := targetConfig.Namespace
	if len(releaseNamespace) == 0 {
		panic("no namespace to deploy to")
	}
	fmt.Printf("releaseNamespace=%s\n", releaseNamespace)

	nexusClient, err := nexus.NewClient(
		*nexusURLFlag,
		*nexusUsernameFlag,
		*nexusPasswordFlag,
		project,
	)
	check(err)
	nexusGroupPrefix := fmt.Sprintf("/%s/%s", repository, gitCommitSHA)

	fmt.Println("copying images...")

	urls, _ := nexusClient.URLs(
		fmt.Sprintf("%s/image-digests", nexusGroupPrefix),
	)

	for _, u := range urls {
		imageStream := strings.TrimSuffix(filepath.Base(u), ".json")
		fmt.Println("copying image", imageStream)
		// TODO: should we also allow external registries? maybe not ...
		srcImageStreamUrl, err := getImageStreamUrl(namespace, imageStream)
		srcRegistryTLSVerify := false
		check(err)
		// TODO: At least for OpenShift image streams, we want to autocreate
		// the destination if it does not exist yet.
		var destImageStreamUrl string
		destRegistryTLSVerify := true
		if len(targetConfig.RegistryHost) > 0 {
			destImageStreamUrl = fmt.Sprintf("%s/%s/%s", targetConfig.RegistryHost, releaseNamespace, imageStream)
			if targetConfig.RegistryTLSVerify != nil {
				destRegistryTLSVerify = *targetConfig.RegistryTLSVerify
			}
		} else {
			disu, err := getImageStreamUrl(releaseNamespace, imageStream)
			check(err)
			destImageStreamUrl = disu
			destRegistryTLSVerify = false
		}
		fmt.Printf("srcRegistry=%s\n", srcImageStreamUrl)
		fmt.Printf("destRegistry=%s\n", destImageStreamUrl)
		// TODO: for QA and PROD we want to ensure that the SHA recorded in Nexus
		// matches the SHA referenced by the Git commit tag.
		stdout, stderr, err := runCmd(
			"skopeo",
			[]string{
				fmt.Sprintf("--src-tls-verify=%v", srcRegistryTLSVerify),
				fmt.Sprintf("--dest-tls-verify=%v", destRegistryTLSVerify),
				fmt.Sprintf("docker://%s:%s", srcImageStreamUrl, gitCommitSHA),
				fmt.Sprintf("docker://%s:%s", destImageStreamUrl, gitCommitSHA),
			},
		)
		if err != nil {
			fmt.Println(err)
			fmt.Println(string(stderr))
		} else {
			fmt.Println(string(stdout))
			fmt.Println(string(stderr))
		}
	}

	fmt.Println("list helm plugins...")

	stdout, stderr, err := runCmd(
		"helm",
		[]string{
			"plugin",
			"list",
		},
	)
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(stderr))
	} else {
		fmt.Println(string(stdout))
		fmt.Println(string(stderr))
	}

	// if child repos exist, collect helm charts for them and place into charts/
	if len(odsConfig.Repositories) > 0 {
		cwd, err := os.Getwd()
		check(err)
		fmt.Println("pulling in helm chart packages from child repositories ...")
		for _, childRepo := range odsConfig.Repositories {
			childCommitSHA, err := getGitCommitSHAInDir(gitRefSpec, ".ods/repositories/"+childRepo.Name)
			check(err)
			// TODO: This should only return one URL - should we enforce this?
			helmChartURLs, err := nexusClient.URLs(
				fmt.Sprintf("/%s/%s/helm-charts", childRepo.Name, childCommitSHA),
			)
			check(err)
			// helm pull
			chartsPath := filepath.Join(*chartDir, "charts")
			err = os.Mkdir(chartsPath, 0644)
			check(err)
			err = os.Chdir(chartsPath)
			check(err)
			for _, helmChartURL := range helmChartURLs {
				nexusClient.Download(helmChartURL)
			}
			err = os.Chdir(cwd)
			check(err)
		}
	}

	fmt.Println("packaging helm chart ...")
	chartVersion, err := getChartVersion(filepath.Join(*chartDir, "Chart.yaml"))
	check(err)
	_, _, err = runCmd(
		"helm",
		[]string{
			"package",
			fmt.Sprintf("--app-version=%s", gitCommitSHA),
			fmt.Sprintf("--version=%s+%s", chartVersion, gitCommitSHA),
			*chartDir,
		},
	)
	check(err)

	helmArchive, err := getHelmArchive(component)
	check(err)

	fmt.Println("uploading helm chart package ...")
	// TODO: check err
	err = nexusClient.Upload(fmt.Sprintf("%s/helm-charts", nexusGroupPrefix), helmArchive)
	if err != nil {
		fmt.Printf("got err: %s", err)
	}

	fmt.Printf("diffing helm release against %s...\n", helmArchive)
	stdout, stderr, err = runCmd(
		"helm",
		[]string{
			"--namespace=" + releaseNamespace,
			"diff",
			"upgrade",
			"--install",
			"--detailed-exitcode",
			"--no-color",
			releaseName,
			helmArchive,
		},
	)

	if err == nil {
		fmt.Println("no diff ...")
		os.Exit(0)
	} else {
		fmt.Println(string(stdout))
		fmt.Println(string(stderr))
	}

	fmt.Printf("upgrading helm release to %s...\n", helmArchive)
	stdout, stderr, err = runCmd(
		"helm",
		[]string{
			"--namespace=" + releaseNamespace,
			"upgrade",
			"--wait",
			"--install",
			releaseName,
			helmArchive,
		},
	)

	check(err)

	fmt.Println(string(stdout))
	fmt.Println(string(stderr))
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func getGitFullRef() (string, error) {
	return getTrimmedFileContent(".ods/git-full-ref")
}

func getGitCommitSHA(refSpec string) (string, error) {
	return getTrimmedFileContent(".ods/git-commit-sha")
}

func getGitCommitSHAInDir(refSpec string, dir string) (string, error) {
	return getTrimmedFileContent(dir + "/.ods/git-commit-sha")
}

func getGitRepository() (string, error) {
	return getTrimmedFileContent(".ods/repository")
}

func runCmd(executable string, args []string) (outBytes, errBytes []byte, err error) {
	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outBytes = stdout.Bytes()
	errBytes = stderr.Bytes()
	return outBytes, errBytes, err
}

func getConfig(filename string) (config.ODS, error) {
	// read ods.yml
	var odsConfig config.ODS
	y, err := ioutil.ReadFile(filename)
	if err != nil {
		return odsConfig, fmt.Errorf("could not read: %w", err)
	}

	err = yaml.Unmarshal(y, &odsConfig)
	if err != nil {
		return odsConfig, fmt.Errorf("could not unmarshal: %w", err)
	}

	return odsConfig, nil
}

func getTarget(odsConfig config.ODS, environment string, target string) (config.Target, error) {
	fmt.Printf("looking for namespace for env=%s, target=%s\n", environment, target)
	var targ config.Target

	var targets []config.Target
	if environment == "dev" {
		targets = odsConfig.Environments.DEV.Targets
	} else {
		return targ, fmt.Errorf("not yet")
	}

	for _, t := range targets {
		if t.Name == target {
			return t, nil
		}
	}

	return targ, fmt.Errorf("no match")
}

type helmChart struct {
	Version string `json:"version"`
}

func getChartVersion(filename string) (string, error) {
	y, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read chart: %w", err)
	}

	var hc helmChart
	err = yaml.Unmarshal(y, &hc)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal: %w", err)
	}
	return hc.Version, nil
}

func getHelmArchive(name string) (string, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tgz") {
			return file.Name(), nil
		}
	}
	return "", fmt.Errorf("did not find archive for %s", name)
}

func getImageStreamUrl(namespace, imageStream string) (string, error) {
	stdout, _, err := runCmd(
		"oc",
		[]string{
			"--namespace=" + namespace,
			"get",
			"is/" + imageStream,
			"-ojsonpath='{.status.dockerImageRepository}'",
		},
	)
	if err != nil {
		return "", err
	}
	return string(stdout), nil

}
