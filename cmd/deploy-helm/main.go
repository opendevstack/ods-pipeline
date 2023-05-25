package main

import (
	"flag"
	"io/fs"
	"os"

	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"k8s.io/client-go/kubernetes"
)

const (
	helmBin                     = "helm"
	kubernetesServiceaccountDir = "/var/run/secrets/kubernetes.io/serviceaccount"
)

type options struct {
	// Name of the Secret resource holding the API user credentials.
	apiCredentialsSecret string
	// API server of the target cluster, including scheme.
	apiServer string
	// Target K8s namespace (or OpenShift project) to deploy into.
	namespace string
	// Hostname of the target registry to push images to.
	registryHost string
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
	// Whether to perform just a diff without any upgrade.
	diffOnly bool
	// Whether to enable debug mode.
	debug bool
}

type deployHelm struct {
	logger logging.LeveledLoggerInterface
	// Name of helm binary.
	helmBin          string
	opts             options
	releaseName      string
	releaseNamespace string
	targetConfig     *targetEnvironment
	imageDigests     []string
	cliValues        []string
	helmArchive      string
	valuesFiles      []string
	clientset        *kubernetes.Clientset
	subrepos         []fs.DirEntry
	ctxt             *pipelinectxt.ODSContext
}

var defaultOptions = options{
	checkoutDir:          ".",
	chartDir:             "./chart",
	ageKeySecretField:    "key.txt",
	certDir:              defaultCertDir(),
	srcRegistryTLSVerify: true,
	debug:                (os.Getenv("DEBUG") == "true"),
}

type targetEnvironment struct {
	APIServer         string
	APIToken          string
	RegistryHost      string
	RegistryTLSVerify *bool
	Namespace         string
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
	flag.StringVar(&opts.apiServer, "api-server", defaultOptions.apiServer, "API server of the target cluster, including scheme")
	flag.StringVar(&opts.apiCredentialsSecret, "api-credentials-secret", defaultOptions.apiCredentialsSecret, "Name of the Secret resource holding the API user credentials")
	flag.StringVar(&opts.registryHost, "registry-host", defaultOptions.registryHost, "Hostname of the target registry to push images to")
	flag.StringVar(&opts.namespace, "namespace", defaultOptions.namespace, "Target K8s namespace (or OpenShift project) to deploy into")
	flag.StringVar(&opts.certDir, "cert-dir", defaultOptions.certDir, "Use certificates at the specified path to access the registry")
	flag.BoolVar(&opts.srcRegistryTLSVerify, "src-registry-tls-verify", defaultOptions.srcRegistryTLSVerify, "TLS verify source registry")
	flag.BoolVar(&opts.diffOnly, "diff-only", defaultOptions.diffOnly, "Whether to perform only a diff")
	flag.BoolVar(&opts.debug, "debug", defaultOptions.debug, "debug mode")
	flag.Parse()

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	err := (&deployHelm{helmBin: helmBin, logger: logger, opts: opts}).runSteps(
		setupContext(),
		skipOnEmptyNamespace(),
		setReleaseTarget(),
		detectSubrepos(),
		listHelmPlugins(),
		packageHelmChartWithSubcharts(),
		collectValuesFiles(),
		importAgeKey(),
		diffHelmRelease(),
		detectImageDigests(),
		copyImagesIntoReleaseNamespace(),
		upgradeHelmRelease(),
	)
	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
}

func defaultCertDir() string {
	if _, err := os.Stat(kubernetesServiceaccountDir); err == nil {
		return kubernetesServiceaccountDir
	}
	return "/etc/containers/certs.d"
}
