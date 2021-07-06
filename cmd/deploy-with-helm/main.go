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
	tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
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

	// Find image artifacts.
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
		for _, f := range files {
			var imageArtifact artifact.Image
			artifactFile := filepath.Join(imageDigestsDir, f.Name())
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
			srcRegistryTLSVerify := false
			// TODO: At least for OpenShift image streams, we want to autocreate
			// the destination if it does not exist yet.
			var destImageURL string
			destRegistryTLSVerify := true
			if len(targetConfig.RegistryHost) > 0 {
				destImageURL = fmt.Sprintf("%s/%s/%s", targetConfig.RegistryHost, releaseNamespace, imageStream)
				if targetConfig.RegistryTLSVerify != nil {
					destRegistryTLSVerify = *targetConfig.RegistryTLSVerify
				}
			} else {
				destImageURL = strings.Replace(imageArtifact.Image, "/"+imageArtifact.Repository+"/", "/"+releaseNamespace+"/", -1)
				destRegistryTLSVerify = false
			}
			fmt.Printf("src=%s\n", srcImageURL)
			fmt.Printf("dest=%s\n", destImageURL)
			// TODO: for QA and PROD we want to ensure that the SHA recorded in Nexus
			// matches the SHA referenced by the Git commit tag.
			skopeoCopyArgs := []string{
				"copy",
				fmt.Sprintf("--src-tls-verify=%v", srcRegistryTLSVerify),
				fmt.Sprintf("--dest-tls-verify=%v", destRegistryTLSVerify),
				fmt.Sprintf("docker://%s", srcImageURL),
				fmt.Sprintf("docker://%s", destImageURL),
			}
			if len(destRegistryToken) > 0 {
				skopeoCopyArgs = append(skopeoCopyArgs, "--dest-registry-token", destRegistryToken)
			}
			stdout, stderr, err := command.Run("skopeo", skopeoCopyArgs)
			if err != nil {
				fmt.Println(string(stderr))
				log.Fatal(err)
			}
			fmt.Println(string(stdout))
		}
	}

	fmt.Println("List Helm plugins...")

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

	fmt.Println("Adding dependencies from subrepos into the charts/ directory ...")
	// Find subcharts
	var subrepos []fs.FileInfo
	subreposDir := ".ods/repos"
	if _, err := os.Stat(subreposDir); os.IsNotExist(err) {
		fmt.Printf("no subrepos in %s\n", subreposDir)
	} else {
		f, err := ioutil.ReadDir(subreposDir)
		if err != nil {
			log.Fatal(err)
		}
		subrepos = f
	}
	chartsDir := filepath.Join(*chartDir, "charts")
	gitCommitSHAs := map[string]interface{}{
		"gitCommitSha": ctxt.GitCommitSHA,
	}
	for _, r := range subrepos {
		subrepo := filepath.Join(subreposDir, r.Name())
		subchart := filepath.Join(subrepo, *chartDir)
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
		helmArchive, err := packageHelmChart(subchart, ctxt.Version, gitCommitSHA)
		if err != nil {
			log.Fatal(err)
		}
		helmArchiveName := filepath.Base(helmArchive)
		fmt.Printf("copying %s into %s\n", helmArchiveName, chartsDir)
		file.Copy(helmArchive, filepath.Join(chartsDir, helmArchiveName))
	}

	fmt.Println("Packaging Helm chart ...")
	helmArchive, err := packageHelmChart(*chartDir, ctxt.Version, ctxt.GitCommitSHA)
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
	valuesFilesCandidates := []string{fmt.Sprintf("values.%s.yaml", targetConfig.Kind)}
	if targetConfig.Kind != targetConfig.Name {
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
	for _, vf := range valuesFiles {
		helmDiffArgs = append(helmDiffArgs, fmt.Sprintf("--values=%s", vf))
	}
	stdout, stderr, err = command.Run(
		"helm", append(helmDiffArgs, releaseName, helmArchive),
	)

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
	for _, vf := range valuesFiles {
		helmUpgradeArgs = append(helmUpgradeArgs, fmt.Sprintf("--values=%s", vf))
	}
	stdout, stderr, err = command.Run(
		"helm", append(helmUpgradeArgs, releaseName, helmArchive),
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
	if len(contextVersion) > 0 {
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

func packageHelmChart(chartDir, ctxtVersion, gitCommitSHA string) (string, error) {
	hc, err := getHelmChart(filepath.Join(chartDir, "Chart.yaml"))
	if err != nil {
		return "", fmt.Errorf("could not read chart: %w", err)
	}
	chartVersion := getChartVersion(ctxtVersion, hc)
	packageVersion := fmt.Sprintf("%s+%s", chartVersion, gitCommitSHA)
	stdout, stderr, err := command.Run(
		"helm",
		[]string{
			"package",
			fmt.Sprintf("--app-version=%s", gitCommitSHA),
			fmt.Sprintf("--version=%s", packageVersion),
			chartDir,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"could not package chart %s. stderr: %s, err: %s", chartDir, string(stderr), err,
		)
	}
	fmt.Println(string(stdout))

	helmArchive := fmt.Sprintf("%s-%s.tgz", hc.Name, packageVersion)
	return helmArchive, nil
}
