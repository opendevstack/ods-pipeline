package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"sigs.k8s.io/yaml"
)

func main() {
	chartDir := flag.String("chart-dir", "", "Chart dir")
	environment := flag.String("environment", "", "environment")
	releaseNameFlag := flag.String("release-name", "", "release-name")
	target := flag.String("target", "", "target")
	flag.Parse()

	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(".")
	if err != nil {
		log.Fatal(err)
	}

	var releaseName string
	if len(*releaseNameFlag) > 0 {
		releaseName = *releaseNameFlag
	} else {
		releaseName = ctxt.Component
	}
	fmt.Printf("releaseName=%s\n", releaseName)

	// read ods.yml
	odsConfig, err := getConfig("ods.yml")
	if err != nil {
		log.Fatal(fmt.Sprintf("err during ods config reading %s", err))
	}
	targetConfig, err := getTarget(odsConfig, *environment, *target)
	if err != nil {
		log.Fatal(fmt.Sprintf("err during namespace extraction %s", err))
	}

	releaseNamespace := targetConfig.Namespace
	if len(releaseNamespace) == 0 {
		log.Fatal("no namespace to deploy to")
	}
	fmt.Printf("releaseNamespace=%s\n", releaseNamespace)

	var files []fs.FileInfo
	imageDigestsDir := ".ods/artifacts/image-digests"
	if _, err := os.Stat(imageDigestsDir); os.IsNotExist(err) {
		fmt.Printf("no image digest in %s\n", imageDigestsDir)
	} else {
		f, err := ioutil.ReadDir(imageDigestsDir)
		if err != nil {
			log.Fatal(err)
		}
		files = f
	}

	for _, f := range files {
		filename := f.Name()
		imageStream := strings.TrimSuffix(filepath.Base(filename), ".json")
		fmt.Println("copying image", imageStream)
		// TODO: should we also allow external registries? maybe not ...
		srcImageStreamUrl, err := getImageStreamUrl(ctxt.Namespace, imageStream)
		srcRegistryTLSVerify := false
		if err != nil {
			log.Fatal(err)
		}
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
			if err != nil {
				log.Fatal(err)
			}
			destImageStreamUrl = disu
			destRegistryTLSVerify = false
		}
		fmt.Printf("srcRegistry=%s\n", srcImageStreamUrl)
		fmt.Printf("destRegistry=%s\n", destImageStreamUrl)
		// TODO: for QA and PROD we want to ensure that the SHA recorded in Nexus
		// matches the SHA referenced by the Git commit tag.
		stdout, stderr, err := command.Run(
			"skopeo",
			[]string{
				fmt.Sprintf("--src-tls-verify=%v", srcRegistryTLSVerify),
				fmt.Sprintf("--dest-tls-verify=%v", destRegistryTLSVerify),
				fmt.Sprintf("docker://%s:%s", srcImageStreamUrl, ctxt.GitCommitSHA),
				fmt.Sprintf("docker://%s:%s", destImageStreamUrl, ctxt.GitCommitSHA),
			},
		)
		if err != nil {
			fmt.Println(string(stderr))
			log.Fatal(err)
		}
		fmt.Println(string(stdout))
	}

	fmt.Println("list helm plugins...")

	stdout, stderr, err := command.Run(
		"helm",
		[]string{
			"plugin",
			"list",
		},
	)
	if err != nil {
		fmt.Println(string(stderr))
		log.Fatal(err)
	}
	fmt.Println(string(stdout))

	// if child repos exist, collect helm charts for them and place into charts/
	// if len(odsConfig.Repositories) > 0 {
	// 	cwd, err := os.Getwd()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("pulling in helm chart packages from child repositories ...")
	// 	for _, childRepo := range odsConfig.Repositories {
	// 		childCommitSHA, err := getGitCommitSHAInDir(".ods/repositories/" + childRepo.Name)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		// TODO: This should only return one URL - should we enforce this?
	// 		helmChartURLs, err := nexusClient.URLs(
	// 			fmt.Sprintf("/%s/%s/helm-charts", childRepo.Name, childCommitSHA),
	// 		)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		// helm pull
	// 		chartsPath := filepath.Join(*chartDir, "charts")
	// 		err = os.Mkdir(chartsPath, 0644)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		err = os.Chdir(chartsPath)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		for _, helmChartURL := range helmChartURLs {
	// 			nexusClient.Download(helmChartURL)
	// 		}
	// 		err = os.Chdir(cwd)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	// }

	fmt.Println("packaging helm chart ...")
	chartVersion, err := getChartVersion(filepath.Join(*chartDir, "Chart.yaml"))
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = command.Run(
		"helm",
		[]string{
			"package",
			fmt.Sprintf("--app-version=%s", ctxt.GitCommitSHA),
			fmt.Sprintf("--version=%s+%s", chartVersion, ctxt.GitCommitSHA),
			*chartDir,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	helmArchive, err := getHelmArchive(ctxt.Component)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("uploading helm chart package ...")
	// // TODO: check err
	// err = nexusClient.Upload(fmt.Sprintf("%s/helm-charts", nexusGroupPrefix), helmArchive)
	// if err != nil {
	// 	fmt.Printf("got err: %s", err)
	// }

	fmt.Printf("diffing helm release against %s...\n", helmArchive)
	stdout, stderr, err = command.Run(
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
	}
	fmt.Println(string(stdout))
	fmt.Println(string(stderr))

	fmt.Printf("upgrading helm release to %s...\n", helmArchive)
	stdout, stderr, err = command.Run(
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
	if err != nil {
		fmt.Println(string(stderr))
		log.Fatal(err)
	}
	fmt.Println(string(stdout))
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
	stdout, _, err := command.Run(
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
