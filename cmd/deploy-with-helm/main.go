package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/file"
	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	tokenFile                   = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	helmBin                     = "helm"
	kubernetesServiceaccountDir = "/var/run/secrets/kubernetes.io/serviceaccount"
)

type options struct {
	chartDir             string
	releaseName          string
	certDir              string
	srcRegistryTLSVerify bool
	debug                bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.chartDir, "chart-dir", "", "Chart dir")
	flag.StringVar(&opts.releaseName, "release-name", "", "release-name")
	flag.StringVar(&opts.certDir, "cert-dir", "/etc/containers/certs.d", "Use certificates at the specified path to access the registry")
	flag.BoolVar(&opts.srcRegistryTLSVerify, "src-registry-tls-verify", true, "TLS verify source registry")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "debug mode")
	flag.Parse()

	checkoutDir := "."

	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}

	if len(ctxt.Environment) == 0 {
		fmt.Println("No environment to deploy to selected. Skipping deployment ...")
		return
	}

	if _, err := os.Stat(kubernetesServiceaccountDir); err == nil {
		opts.certDir = kubernetesServiceaccountDir
	}
	if opts.debug {
		directory.ListFiles(opts.certDir)
	}

	var releaseName string
	if len(opts.releaseName) > 0 {
		releaseName = opts.releaseName
	} else {
		releaseName = ctxt.Component
	}
	fmt.Printf("releaseName=%s\n", releaseName)

	// read ods.yml
	odsConfig, err := config.ReadFromDir(checkoutDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("err during ods config reading: %s", err))
	}
	targetConfig, err := getTargetEnvironment(odsConfig, ctxt.Environment)
	if err != nil {
		log.Fatal(fmt.Sprintf("err during namespace extraction: %s", err))
	}

	releaseNamespace := targetConfig.Namespace
	if len(releaseNamespace) == 0 {
		releaseNamespace = fmt.Sprintf("%s-%s", ctxt.Project, targetConfig.Name)
	}
	fmt.Printf("releaseNamespace=%s\n", releaseNamespace)

	// Find subrepos
	var subrepos []fs.FileInfo
	if _, err := os.Stat(pipelinectxt.SubreposPath); err == nil {
		f, err := ioutil.ReadDir(pipelinectxt.SubreposPath)
		if err != nil {
			log.Fatal(err)
		}
		subrepos = f
	}

	// Find image artifacts.
	var files []string
	id, err := collectImageDigests(pipelinectxt.ImageDigestsPath)
	if err != nil {
		log.Fatal(err)
	}
	files = append(files, id...)
	for _, s := range subrepos {
		subrepoImageDigestsPath := filepath.Join(pipelinectxt.SubreposPath, s.Name(), pipelinectxt.ImageDigestsPath)
		id, err := collectImageDigests(subrepoImageDigestsPath)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, id...)
	}

	// Copy images into release namespace if there are any image artifacts.
	if len(files) > 0 {
		clientset, err := k.NewInClusterClientset()
		if err != nil {
			log.Fatalf("could not create Kubernetes client: %s", err)
		}
		// Get destination registry token from secret or file in pod.
		var destRegistryToken string
		if len(targetConfig.SecretRef) > 0 {
			token, err := tokenFromSecret(clientset, releaseNamespace, targetConfig.SecretRef)
			if err != nil {
				log.Fatalf("could not get token from secret %s: %s", targetConfig.SecretRef, err)
			}
			destRegistryToken = token
		} else {
			token, err := getTrimmedFileContent(tokenFile)
			if err != nil {
				log.Fatalf("could not get token from file %s: %s", tokenFile, err)
			}
			destRegistryToken = token
		}

		fmt.Println("Copying images into release namespace ...")
		for _, artifactFile := range files {
			var imageArtifact artifact.Image
			artifactContent, err := ioutil.ReadFile(artifactFile)
			if err != nil {
				log.Fatalf("could not read image artifact file %s: %s", artifactFile, err)
			}
			err = json.Unmarshal(artifactContent, &imageArtifact)
			if err != nil {
				log.Fatalf(
					"could not unmarshal image artifact file %s: %s.\nFile content:\n%s",
					artifactFile, err, string(artifactContent),
				)
			}
			imageStream := imageArtifact.Name
			fmt.Println("copying image", imageStream)
			srcImageURL := imageArtifact.Image
			// TODO: At least for OpenShift image streams, we want to autocreate
			// the destination if it does not exist yet.
			var destImageURL string
			// If the source registry should be TLS verified, the destination
			// should be verified by default as well.
			destRegistryTLSVerify := opts.srcRegistryTLSVerify
			if len(targetConfig.RegistryHost) > 0 {
				destImageURL = fmt.Sprintf("%s/%s/%s", targetConfig.RegistryHost, releaseNamespace, imageStream)
				if targetConfig.RegistryTLSVerify != nil {
					destRegistryTLSVerify = *targetConfig.RegistryTLSVerify
				}
			} else {
				destImageURL = strings.Replace(imageArtifact.Image, "/"+imageArtifact.Repository+"/", "/"+releaseNamespace+"/", -1)
			}
			fmt.Printf("src=%s\n", srcImageURL)
			fmt.Printf("dest=%s\n", destImageURL)
			// TODO: for QA and PROD we want to ensure that the SHA recorded in Nexus
			// matches the SHA referenced by the Git commit tag.
			skopeoCopyArgs := []string{
				"copy",
				fmt.Sprintf("--src-tls-verify=%v", opts.srcRegistryTLSVerify),
				fmt.Sprintf("--dest-tls-verify=%v", destRegistryTLSVerify),
			}
			if opts.srcRegistryTLSVerify {
				skopeoCopyArgs = append(skopeoCopyArgs, fmt.Sprintf("--src-cert-dir=%v", opts.certDir))
			}
			if destRegistryTLSVerify {
				skopeoCopyArgs = append(skopeoCopyArgs, fmt.Sprintf("--dest-cert-dir=%v", opts.certDir))
			}
			if len(destRegistryToken) > 0 {
				skopeoCopyArgs = append(skopeoCopyArgs, "--dest-registry-token", destRegistryToken)
			}
			if opts.debug {
				skopeoCopyArgs = append(skopeoCopyArgs, "--debug")
			}
			stdout, stderr, err := command.Run(
				"skopeo", append(
					skopeoCopyArgs,
					fmt.Sprintf("docker://%s", srcImageURL),
					fmt.Sprintf("docker://%s", destImageURL),
				),
			)
			if err != nil {
				fmt.Println(string(stderr))
				log.Fatal(err)
			}
			fmt.Println(string(stdout))
		}
	}

	fmt.Println("List Helm plugins...")
	helmPluginArgs := []string{"plugin", "list"}
	if opts.debug {
		helmPluginArgs = append(helmPluginArgs, "--debug")
	}
	stdout, stderr, err := command.Run(helmBin, helmPluginArgs)
	if err != nil {
		fmt.Println(string(stderr))
		log.Fatal(err)
	}
	fmt.Println(string(stdout))

	fmt.Println("Adding dependencies from subrepos into the charts/ directory ...")
	// Find subcharts
	chartsDir := filepath.Join(opts.chartDir, "charts")
	if _, err := os.Stat(chartsDir); os.IsNotExist(err) {
		err = os.Mkdir(chartsDir, 0755)
		if err != nil {
			log.Fatalf("could not create %s: %s", chartsDir, err)
		}
	}
	gitCommitSHAs := map[string]interface{}{
		"gitCommitSha": ctxt.GitCommitSHA,
	}
	for _, r := range subrepos {
		subrepo := filepath.Join(pipelinectxt.SubreposPath, r.Name())
		subchart := filepath.Join(subrepo, opts.chartDir)
		if _, err := os.Stat(subchart); os.IsNotExist(err) {
			fmt.Printf("no chart in %s\n", r.Name())
			continue
		}
		gitCommitSHA, err := getTrimmedFileContent(filepath.Join(subrepo, ".ods", "git-commit-sha"))
		if err != nil {
			log.Fatal(err)
		}
		hc, err := getHelmChart(filepath.Join(subchart, "Chart.yaml"))
		if err != nil {
			log.Fatal(err)
		}
		gitCommitSHAs[hc.Name] = map[string]string{
			"gitCommitSha": gitCommitSHA,
		}
		helmArchive, err := packageHelmChart(subchart, ctxt.Version, gitCommitSHA, opts.debug)
		if err != nil {
			log.Fatal(err)
		}
		helmArchiveName := filepath.Base(helmArchive)
		fmt.Printf("copying %s into %s\n", helmArchiveName, chartsDir)
		err = file.Copy(helmArchive, filepath.Join(chartsDir, helmArchiveName))
		if err != nil {
			log.Fatal(err)
		}
	}

	subcharts, err := ioutil.ReadDir(chartsDir)
	if err != nil {
		log.Fatal(err)
	}
	if len(subcharts) > 0 {
		fmt.Printf("Contents of %s:\n", chartsDir)
		for _, sc := range subcharts {
			fmt.Println(sc.Name())
		}
	}

	fmt.Println("Packaging Helm chart ...")
	helmArchive, err := packageHelmChart(opts.chartDir, ctxt.Version, ctxt.GitCommitSHA, opts.debug)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Collecting Helm values files ...")
	generatedValuesFilename := "values.generated.yaml"
	out, err := yaml.Marshal(gitCommitSHAs)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(generatedValuesFilename, out, 0644)
	if err != nil {
		log.Fatal(err)
	}
	valuesFiles := []string{}
	valuesFilesCandidates := []string{fmt.Sprintf("values.%s.yaml", targetConfig.Stage)}
	if string(targetConfig.Stage) != targetConfig.Name {
		valuesFilesCandidates = append(valuesFilesCandidates, fmt.Sprintf("values.%s.yaml", targetConfig.Name))
	}
	valuesFilesCandidates = append(valuesFilesCandidates, generatedValuesFilename)
	for _, vfc := range valuesFilesCandidates {
		if _, err := os.Stat(vfc); os.IsNotExist(err) {
			fmt.Printf("%s is not present, skipping.\n", vfc)
		} else {
			fmt.Printf("%s is present, adding.\n", vfc)
			valuesFiles = append(valuesFiles, vfc)
		}
	}

	fmt.Printf("Diffing Helm release against %s...\n", helmArchive)
	helmDiffArgs := []string{
		"--namespace=" + releaseNamespace,
		"diff",
		"upgrade",
		"--install",
		"--detailed-exitcode",
		"--no-color",
	}
	if opts.debug {
		helmDiffArgs = append(helmDiffArgs, "--debug")
	}
	for _, vf := range valuesFiles {
		helmDiffArgs = append(helmDiffArgs, fmt.Sprintf("--values=%s", vf))
	}
	helmDiffArgs = append(helmDiffArgs, releaseName, helmArchive)
	fmt.Println(helmBin, strings.Join(helmDiffArgs, " "))
	stdout, stderr, err = command.Run(helmBin, helmDiffArgs)

	if err == nil {
		fmt.Println("no diff ...")
		os.Exit(0)
	}
	fmt.Println(string(stdout))
	fmt.Println(string(stderr))

	fmt.Printf("Upgrading Helm release to %s...\n", helmArchive)
	helmUpgradeArgs := []string{
		"--namespace=" + releaseNamespace,
		"upgrade",
		"--wait",
		"--install",
	}
	if opts.debug {
		helmUpgradeArgs = append(helmUpgradeArgs, "--debug")
	}
	for _, vf := range valuesFiles {
		helmUpgradeArgs = append(helmUpgradeArgs, fmt.Sprintf("--values=%s", vf))
	}
	helmUpgradeArgs = append(helmUpgradeArgs, releaseName, helmArchive)
	fmt.Println(helmBin, strings.Join(helmUpgradeArgs, " "))
	stdout, stderr, err = command.Run(helmBin, helmUpgradeArgs)
	if err != nil {
		fmt.Println(string(stderr))
		log.Fatal(err)
	}
	fmt.Println(string(stdout))
}

func getTargetEnvironment(odsConfig *config.ODS, environment string) (*config.Environment, error) {
	var envs []string
	for _, e := range odsConfig.Environments {
		if e.Name == environment {
			return &e, nil
		}
		envs = append(envs, e.Name)
	}

	return nil, fmt.Errorf("no environment matched '%s', have: %s", environment, strings.Join(envs, ", "))
}

type helmChart struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func getHelmChart(filename string) (*helmChart, error) {
	y, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read chart: %w", err)
	}

	var hc *helmChart
	err = yaml.Unmarshal(y, &hc)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
	}
	return hc, nil
}

func getChartVersion(contextVersion string, hc *helmChart) string {
	if len(contextVersion) > 0 && contextVersion != pipelinectxt.WIP {
		return contextVersion
	}
	return hc.Version
}

func tokenFromSecret(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data["token"]), nil
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func packageHelmChart(chartDir, ctxtVersion, gitCommitSHA string, debug bool) (string, error) {
	hc, err := getHelmChart(filepath.Join(chartDir, "Chart.yaml"))
	if err != nil {
		return "", fmt.Errorf("could not read chart: %w", err)
	}
	chartVersion := getChartVersion(ctxtVersion, hc)
	packageVersion := fmt.Sprintf("%s+%s", chartVersion, gitCommitSHA)
	helmPackageArgs := []string{
		"package",
		fmt.Sprintf("--app-version=%s", gitCommitSHA),
		fmt.Sprintf("--version=%s", packageVersion),
	}
	if debug {
		helmPackageArgs = append(helmPackageArgs, "--debug")
	}
	stdout, stderr, err := command.Run(helmBin, append(helmPackageArgs, chartDir))
	if err != nil {
		return "", fmt.Errorf(
			"could not package chart %s. stderr: %s, err: %s", chartDir, string(stderr), err,
		)
	}
	fmt.Println(string(stdout))

	helmArchive := fmt.Sprintf("%s-%s.tgz", hc.Name, packageVersion)
	return helmArchive, nil
}

func collectImageDigests(imageDigestsDir string) ([]string, error) {
	var files []string
	if _, err := os.Stat(imageDigestsDir); err == nil {
		f, err := ioutil.ReadDir(imageDigestsDir)
		if err != nil {
			return files, fmt.Errorf("could not read image digests dir: %w", err)
		}
		for _, fi := range f {
			files = append(files, filepath.Join(imageDigestsDir, fi.Name()))
		}
	}
	return files, nil
}
