package main

import (
	"bytes"
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

	"github.com/google/shlex"
	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/file"
	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	tokenFile                   = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	helmBin                     = "helm"
	kubernetesServiceaccountDir = "/var/run/secrets/kubernetes.io/serviceaccount"
	// file path where to internally store the age-key-secret openshift secret content,
	// required by helm secrets plugin.
	ageKeyFilePath = "./key.txt"
)

// confusingHelmDiffMessage is the message Helm prints when helm-diff is
// configured to exit with a non-zero exit code when drift is detected,
// provided helm-secrets is also installed but no secrets are decrypted.
const confusingHelmDiffMessage = `Error: identified at least one change, exiting with non-zero exit code (detailed-exitcode parameter enabled)
Error: plugin "diff" exited with error
Error: plugin "secrets" exited with error`

// confusingHelmDiffDecryptionMessage is like confusingHelmDiffMessage but occurs
// when helm-secrets is required to decrypt secrets to perform the diff.
const confusingHelmDiffDecryptionMessage = `Error: identified at least one change, exiting with non-zero exit code (detailed-exitcode parameter enabled)
Error: plugin "diff" exited with error

/usr/local/helm/plugins/helm-secrets/scripts/commands/helm.sh: line 34: xargs: command not found
Error: plugin "secrets" exited with error`

type options struct {
	// Location of Helm chart directory.
	chartDir string
	// Name of Helm release.
	releaseName string
	// Flags to pass to `helm diff upgrade` (in addition to default ones).
	diffFlags string
	// Flags to pass to `helm upgrade`.
	upgradeFlags string
	// Name of K8s secret holding the age key.
	ageKeySecret string
	// Field name within the K8s secret holding the age key.
	ageKeySecretField string
	// Location of the certificate directory.
	certDir string
	// Whether to TLS verify the source image registry.
	srcRegistryTLSVerify bool
	// Whether to enable debug mode.
	debug bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.chartDir, "chart-dir", "", "Chart dir")
	flag.StringVar(&opts.releaseName, "release-name", "", "Name of Helm release")
	flag.StringVar(&opts.diffFlags, "diff-flags", "", "Flags to pass to `helm diff upgrade` (in addition to default ones)")
	flag.StringVar(&opts.upgradeFlags, "upgrade-flags", "", "Flags to pass to `helm upgrade`")
	flag.StringVar(&opts.ageKeySecret, "age-key-secret", "", "Name of the secret containing the age key to use for helm-secrets")
	flag.StringVar(&opts.ageKeySecretField, "age-key-secret-field", "key.txt", "Name of the field in the secret holding the age private key")
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

	// read ods.y(a)ml
	odsConfig, err := config.ReadFromDir(checkoutDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("err during ods config reading: %s", err))
	}
	targetConfig, err := odsConfig.Environment(ctxt.Environment)
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

	clientset, err := k.NewInClusterClientset()
	if err != nil {
		log.Fatalf("could not create Kubernetes client: %s", err)
	}

	if targetConfig.APIServer != "" {
		token, err := tokenFromSecret(clientset, ctxt.Namespace, targetConfig.APICredentialsSecret)
		if err != nil {
			log.Fatalf("could not get token from secret %s: %s", targetConfig.APICredentialsSecret, err)
		}
		targetConfig.APIToken = token
	}

	// Copy images into release namespace if there are any image artifacts.
	if len(files) > 0 {
		// Get destination registry token from secret or file in pod.
		var destRegistryToken string
		if targetConfig.APIToken != "" {
			destRegistryToken = targetConfig.APIToken
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
			fmt.Printf("Copying image %s ...\n", imageStream)
			srcImageURL := imageArtifact.Image
			var destImageURL string
			// If the source registry should be TLS verified, the destination
			// should be verified by default as well.
			destRegistryTLSVerify := opts.srcRegistryTLSVerify
			srcRegistryTLSVerify := opts.srcRegistryTLSVerify
			// TLS verification of the KinD registry is not possible at the moment as
			// requests error out with "server gave HTTP response to HTTPS client".
			if strings.HasPrefix(imageArtifact.Registry, "kind-registry.kind") {
				srcRegistryTLSVerify = false
			}
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
				fmt.Sprintf("--src-tls-verify=%v", srcRegistryTLSVerify),
				fmt.Sprintf("--dest-tls-verify=%v", destRegistryTLSVerify),
			}
			if srcRegistryTLSVerify {
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

	// Collect values to be set via the CLI.
	cliValues := []string{
		fmt.Sprintf("--set=image.tag=%s", ctxt.GitCommitSHA),
	}

	fmt.Println("Adding dependencies from subrepos into the charts/ directory ...")
	// Find subcharts
	chartsDir := filepath.Join(opts.chartDir, "charts")
	if _, err := os.Stat(chartsDir); os.IsNotExist(err) {
		err = os.Mkdir(chartsDir, 0755)
		if err != nil {
			log.Fatalf("could not create %s: %s", chartsDir, err)
		}
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
		cliValues = append(cliValues, fmt.Sprintf("--set=%s.image.tag=%s", hc.Name, gitCommitSHA))
		if releaseName == ctxt.Component {
			cliValues = append(cliValues, fmt.Sprintf("--set=%s.fullnameOverride=%s", hc.Name, hc.Name))
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
	valuesFiles := []string{}
	valuesFilesCandidates := []string{
		fmt.Sprintf("%s/secrets.yaml", opts.chartDir), // equivalent values.yaml is added automatically by Helm
		fmt.Sprintf("%s/values.%s.yaml", opts.chartDir, targetConfig.Stage),
		fmt.Sprintf("%s/secrets.%s.yaml", opts.chartDir, targetConfig.Stage),
	}
	if string(targetConfig.Stage) != targetConfig.Name {
		valuesFilesCandidates = append(
			valuesFilesCandidates,
			fmt.Sprintf("%s/values.%s.yaml", opts.chartDir, targetConfig.Name),
			fmt.Sprintf("%s/secrets.%s.yaml", opts.chartDir, targetConfig.Name),
		)
	}
	for _, vfc := range valuesFilesCandidates {
		if _, err := os.Stat(vfc); os.IsNotExist(err) {
			fmt.Printf("%s is not present, skipping.\n", vfc)
		} else {
			fmt.Printf("%s is present, adding.\n", vfc)
			valuesFiles = append(valuesFiles, vfc)
		}
	}

	if len(opts.ageKeySecret) == 0 {
		fmt.Println("Skipping import of age key for helm-secrets as parameter is not set ...")
	} else {
		fmt.Println("Storing age key for helm-secrets ...")
		secret, err := clientset.CoreV1().Secrets(ctxt.Namespace).Get(
			context.TODO(), opts.ageKeySecret, metav1.GetOptions{},
		)
		if err != nil {
			fmt.Printf("No secret %s found, skipping.\n", opts.ageKeySecret)
		} else {
			stderr, err = storeAgeKey(secret, opts.ageKeySecretField)
			if err != nil {
				fmt.Println(string(stderr))
				log.Fatal(err)
			}
			fmt.Printf("Age key secret %s stored.\n", opts.ageKeySecret)
		}
	}

	fmt.Printf("Diffing Helm release against %s...\n", helmArchive)
	helmDiffArgs := []string{
		"--namespace=" + releaseNamespace,
		"secrets",
		"diff",
		"upgrade",
		"--detailed-exitcode",
		"--no-color",
	}
	helmDiffFlags, err := shlex.Split(opts.diffFlags)
	if err != nil {
		log.Fatalf("could not parse diff flags (%s): %s", opts.diffFlags, err)
	}
	helmDiffArgs = append(helmDiffArgs, helmDiffFlags...)
	for _, vf := range valuesFiles {
		helmDiffArgs = append(helmDiffArgs, fmt.Sprintf("--values=%s", vf))
	}
	helmDiffArgs = append(helmDiffArgs, cliValues...)
	helmDiffArgs = append(helmDiffArgs, releaseName, helmArchive)
	stdout, stderr, err = runHelmCmd(helmDiffArgs, targetConfig, opts.debug)
	if err == nil {
		fmt.Println("no diff ...")
		os.Exit(0)
	}
	// Replace confusing stderr messages while still printing stderr in general
	// to surface any other issues that might be logged there.
	diffStderr := replaceConfusingHelmLogMessages(stderr)
	fmt.Println(string(stdout))
	fmt.Println(string(diffStderr))
	err = writeDeploymentArtifact(stdout, "diff", opts.chartDir, targetConfig.Name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Upgrading Helm release to %s...\n", helmArchive)
	helmUpgradeArgs := []string{
		"--namespace=" + releaseNamespace,
		"secrets",
		"upgrade",
	}
	helmUpgradeFlags, err := shlex.Split(opts.upgradeFlags)
	if err != nil {
		log.Fatalf("could not parse upgrade flags (%s): %s", opts.upgradeFlags, err)
	}
	helmUpgradeArgs = append(helmUpgradeArgs, helmUpgradeFlags...)
	for _, vf := range valuesFiles {
		helmUpgradeArgs = append(helmUpgradeArgs, fmt.Sprintf("--values=%s", vf))
	}
	helmUpgradeArgs = append(helmUpgradeArgs, cliValues...)
	helmUpgradeArgs = append(helmUpgradeArgs, releaseName, helmArchive)
	stdout, stderr, err = runHelmCmd(helmUpgradeArgs, targetConfig, opts.debug)
	if err != nil {
		fmt.Println(string(stderr))
		log.Fatal(err)
	}
	fmt.Println(string(stdout))
	err = writeDeploymentArtifact(stdout, "release", opts.chartDir, targetConfig.Name)
	if err != nil {
		log.Fatal(err)
	}
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

// replaceConfusingHelmLogMessages replaces confusing stderr messages with
// easier to understand ones. The tests in ods-deploy-helm_test.go ensure this
// still works when Helm is upgraded.
func replaceConfusingHelmLogMessages(logs []byte) []byte {
	cleaned := bytes.Replace(
		logs,
		[]byte(confusingHelmDiffMessage),
		[]byte("plugin \"diff\" identified at least one change"),
		1,
	)
	cleaned = bytes.Replace(
		cleaned,
		[]byte(confusingHelmDiffDecryptionMessage),
		[]byte("plugin \"diff\" identified at least one change"),
		1,
	)
	return cleaned
}

func artifactFilename(filename, chartDir, targetEnv string) string {
	trimmedChartDir := strings.TrimPrefix(chartDir, "./")
	if trimmedChartDir != "chart" {
		filename = fmt.Sprintf("%s-%s", strings.Replace(trimmedChartDir, "/", "-", -1), filename)
	}
	return fmt.Sprintf("%s-%s", filename, targetEnv)
}

func writeDeploymentArtifact(content []byte, filename, chartDir, targetEnv string) error {
	err := os.MkdirAll(pipelinectxt.DeploymentsPath, 0755)
	if err != nil {
		return err
	}
	f := artifactFilename(filename, chartDir, targetEnv) + ".txt"
	return ioutil.WriteFile(filepath.Join(pipelinectxt.DeploymentsPath, f), content, 0644)
}

func storeAgeKey(secret *corev1.Secret, ageKeySecretField string) (errBytes []byte, err error) {
	file, err := os.Create(ageKeyFilePath)
	if err != nil {
		return errBytes, err
	}
	defer file.Close()
	_, err = file.Write(secret.Data[ageKeySecretField])
	if err != nil {
		return errBytes, err
	}
	return errBytes, err
}

func runHelmCmd(args []string, targetConfig *config.Environment, debug bool) (outBytes, errBytes []byte, err error) {
	if debug {
		args = append([]string{"--debug"}, args...)
	}
	printableArgs := args
	if targetConfig.APIServer != "" {
		printableArgs = append(
			[]string{
				fmt.Sprintf("--kube-apiserver=%s", targetConfig.APIServer),
				"--kube-token=***",
			},
			args...,
		)
		args = append(
			[]string{
				fmt.Sprintf("--kube-apiserver=%s", targetConfig.APIServer),
				fmt.Sprintf("--kube-token=%s", targetConfig.APIToken),
			},
			args...,
		)
	}
	fmt.Println(helmBin, strings.Join(printableArgs, " "))

	var extraEnvs = []string{fmt.Sprintf("SOPS_AGE_KEY_FILE=%s", ageKeyFilePath)}
	return command.RunWithExtraEnvs(helmBin, args, extraEnvs)
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
