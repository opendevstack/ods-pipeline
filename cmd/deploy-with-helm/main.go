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

	// Get destination registry token if there are any image artifacts.
	var destRegistryToken string
	if len(files) > 0 {
		clientset, err := k.NewInClusterClientset()
		if err != nil {
			log.Fatalf("could not create Kubernetes client: %s", err)
		}
		if len(targetConfig.SecretRef) > 0 {
			token, err := tokenFromSecret(clientset, releaseNamespace, targetConfig.SecretRef)
			if err != nil {
				log.Fatalf("could not get token from secret %s: %s", targetConfig.SecretRef, err)
			}
			destRegistryToken = token
		} else {
			token, err := tokenFromFile(tokenFile)
			if err != nil {
				log.Fatalf("could not get token from file %s: %s", tokenFile, err)
			}
			destRegistryToken = token
		}
	}

	// Copy images into release namespace.
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

func tokenFromSecret(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data["token"]), nil
}

func tokenFromFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
