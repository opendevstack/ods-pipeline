package main

import (
	"context"
	"encoding/base64"
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
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

type ImageDigest struct {
	Image      string `json:"image"`
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	Digest     string `json:"digest"`
}

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

	destCreds := map[string]dockerConfig{}
	if len(files) > 0 {
		clientset, err := k8sClient()
		if err != nil {
			log.Fatal(err)
		}
		dc, err := saDockercfgs(clientset, releaseNamespace, "builder")
		if err != nil {
			log.Fatal(err)
		}
		destCreds = dc
	}

	for _, f := range files {
		var id ImageDigest
		idContent, err := ioutil.ReadFile(filepath.Join(imageDigestsDir, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(idContent, &id)
		if err != nil {
			log.Fatal(err)
		}
		imageStream := id.Name
		fmt.Println("copying image", imageStream)
		srcImageURL := id.Image
		srcRegistryTLSVerify := false
		// TODO: At least for OpenShift image streams, we want to autocreate
		// the destination if it does not exist yet.
		var destImageURL string
		var destRegistry string
		destRegistryTLSVerify := true
		if len(targetConfig.RegistryHost) > 0 {
			destRegistry = targetConfig.RegistryHost
			destImageURL = fmt.Sprintf("%s/%s/%s", destRegistry, releaseNamespace, imageStream)
			if targetConfig.RegistryTLSVerify != nil {
				destRegistryTLSVerify = *targetConfig.RegistryTLSVerify
			}
		} else {
			destRegistry = id.Registry
			destImageURL = strings.Replace(id.Image, "/"+id.Repository+"/", "/"+releaseNamespace+"/", -1)
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
		if v, ok := destCreds[destRegistry]; ok {
			skopeoCopyArgs = append(skopeoCopyArgs, "--dest-creds", "builder:"+v.Password)
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

func k8sClient() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	return kubernetes.NewForConfig(config)
}

type dockerConfig struct {
	Auth     string `json:"auth"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func saDockercfgs(clientset *kubernetes.Clientset, namespace, serviceaccount string) (map[string]dockerConfig, error) {
	cfg := map[string]dockerConfig{}
	builderServiceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), serviceaccount, metav1.GetOptions{})
	if err != nil {
		return cfg, err
	}
	dockercfgSecretPrefix := serviceaccount + "-dockercfg-"
	for _, s := range builderServiceAccount.Secrets {
		if strings.HasPrefix(s.Name, dockercfgSecretPrefix) {
			builderDockercfgSecret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), s.Name, metav1.GetOptions{})
			if err != nil {
				return cfg, err
			}
			var decoded []byte
			fmt.Println(builderDockercfgSecret.Data)
			fmt.Println(builderDockercfgSecret.Data[".dockercfg"])
			_, err = base64.StdEncoding.Decode(decoded, builderDockercfgSecret.Data[".dockercfg"])
			if err != nil {
				return cfg, err
			}

			err = json.Unmarshal(decoded, &cfg)
			if err != nil {
				return cfg, err
			}

			return cfg, nil
		}
	}
	return cfg, fmt.Errorf("did not find secrets prefixed with %s", dockercfgSecretPrefix)
}
