package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	tokenFile                   = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	helmBin                     = "helm"
	kubernetesServiceaccountDir = "/var/run/secrets/kubernetes.io/serviceaccount"
	// file path where to internally store the age-key-secret openshift secret content,
	// required by helm secrets plugin.
	ageKeyFilePath = "./key.txt"
)

type options struct {
	// Location of checkout directory.
	checkoutDir string
	// Location of Helm chart directory.
	chartDir string
	// Name of Helm release.
	releaseName string
	// Flags to pass to `helm diff upgrade` (in addition to default ones and upgrade flags).
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

type releaseTarget struct {
	name      string
	namespace string
	config    *config.Environment
}

type deployHelm struct {
	logger        logging.LeveledLoggerInterface
	helmBin       string
	opts          options
	releaseTarget *releaseTarget
	imageDigests  []string
	cliValues     []string
	helmArchive   string
	valuesFiles   []string
	clientset     *kubernetes.Clientset
	subrepos      []fs.DirEntry
	ctxt          *pipelinectxt.ODSContext
}

type skipFollowingSteps struct {
	msg string
}

func (e *skipFollowingSteps) Error() string {
	return e.msg
}

func (d *deployHelm) RunSteps(steps ...DeployStep) error {
	var skip *skipFollowingSteps
	var err error
	for _, step := range steps {
		d, err = step(d)
		if err != nil {
			if errors.As(err, &skip) {
				d.logger.Infof(err.Error())
				return nil
			}
			return err
		}
	}
	return nil
}

var defaultOptions = options{
	checkoutDir:       ".",
	chartDir:          "./chart",
	releaseName:       "",
	diffFlags:         "",
	upgradeFlags:      "",
	ageKeySecret:      "",
	ageKeySecretField: "key.txt",
	certDir: func() string {
		if _, err := os.Stat(kubernetesServiceaccountDir); err == nil {
			return kubernetesServiceaccountDir
		}
		return "/etc/containers/certs.d"
	}(),
	srcRegistryTLSVerify: true,
	debug:                (os.Getenv("DEBUG") == "true"),
}

func main() {
	opts := options{}
	flag.StringVar(&opts.checkoutDir, "checkout-dir", defaultOptions.checkoutDir, "Checkout dir")
	flag.StringVar(&opts.chartDir, "chart-dir", defaultOptions.chartDir, "Chart dir")
	flag.StringVar(&opts.releaseName, "release-name", defaultOptions.releaseName, "Name of Helm release")
	flag.StringVar(&opts.diffFlags, "diff-flags", defaultOptions.diffFlags, "Flags to pass to `helm diff upgrade` (in addition to default ones and upgrade flags)")
	flag.StringVar(&opts.upgradeFlags, "upgrade-flags", defaultOptions.upgradeFlags, "Flags to pass to `helm upgrade`")
	flag.StringVar(&opts.ageKeySecret, "age-key-secret", defaultOptions.ageKeySecret, "Name of the secret containing the age key to use for helm-secrets")
	flag.StringVar(&opts.ageKeySecretField, "age-key-secret-field", defaultOptions.ageKeySecretField, "Name of the field in the secret holding the age private key")
	flag.StringVar(&opts.certDir, "cert-dir", defaultOptions.certDir, "Use certificates at the specified path to access the registry")
	flag.BoolVar(&opts.srcRegistryTLSVerify, "src-registry-tls-verify", defaultOptions.srcRegistryTLSVerify, "TLS verify source registry")
	flag.BoolVar(&opts.debug, "debug", defaultOptions.debug, "debug mode")
	flag.Parse()

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	err := (&deployHelm{helmBin: helmBin, logger: logger, opts: opts}).RunSteps(
		setupContext(),
		skipOnEmptyEnv(),
		determineReleaseTarget(),
		determineSubrepos(),
		determineImageDigests(),
		copyImagesIntoReleaseNamespace(),
		listHelmPlugins(),
		packageHelmChartWithSubcharts(),
		collectValuesFiles(),
		importAgeKey(),
		diffHelmRelease(),
		upgradeHelmRelease(),
	)
	if err != nil {
		log.Fatal(err)
	}
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
	return os.WriteFile(filepath.Join(pipelinectxt.DeploymentsPath, f), content, 0644)
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

func tokenFromSecret(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data["token"]), nil
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func collectImageDigests(imageDigestsDir string) ([]string, error) {
	var files []string
	if _, err := os.Stat(imageDigestsDir); err == nil {
		f, err := os.ReadDir(imageDigestsDir)
		if err != nil {
			return files, fmt.Errorf("could not read image digests dir: %w", err)
		}
		for _, fi := range f {
			files = append(files, filepath.Join(imageDigestsDir, fi.Name()))
		}
	}
	return files, nil
}

func getImageDestURL(registryHost, releaseNamespace string, imageArtifact artifact.Image) string {
	if registryHost != "" {
		return fmt.Sprintf("%s/%s/%s:%s", registryHost, releaseNamespace, imageArtifact.Name, imageArtifact.Tag)
	} else {
		return strings.Replace(imageArtifact.Image, "/"+imageArtifact.Repository+"/", "/"+releaseNamespace+"/", -1)
	}
}
