package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

const (
	kubernetesServiceaccountDir = "/var/run/secrets/kubernetes.io/serviceaccount"
)

type options struct {
	checkoutDir           string
	bitbucketAccessToken  string
	bitbucketURL          string
	aquaUsername          string
	aquaPassword          string
	aquaURL               string
	aquaRegistry          string
	imageStream           string
	registry              string
	certDir               string
	imageNamespace        string
	tlsVerify             bool
	storageDriver         string
	format                string
	dockerfile            string
	contextDir            string
	nexusURL              string
	nexusUsername         string
	nexusPassword         string
	buildahBuildExtraArgs string
	buildahPushExtraArgs  string
	trivySBOMExtraArgs    string
	aquasecGate           bool
	debug                 bool
	sbomFormat            string
}

type packageImage struct {
	logger logging.LeveledLoggerInterface
	opts   options
	ctxt   *pipelinectxt.ODSContext
	image  artifact.Image
}

var defaultOptions = options{
	checkoutDir:           ".",
	bitbucketAccessToken:  os.Getenv("BITBUCKET_ACCESS_TOKEN"),
	bitbucketURL:          os.Getenv("BITBUCKET_URL"),
	aquaUsername:          os.Getenv("AQUA_USERNAME"),
	aquaPassword:          os.Getenv("AQUA_PASSWORD"),
	aquaURL:               os.Getenv("AQUA_URL"),
	aquaRegistry:          os.Getenv("AQUA_REGISTRY"),
	imageStream:           "",
	registry:              "image-registry.openshift-image-registry.svc:5000",
	certDir:               defaultCertDir(),
	imageNamespace:        "",
	tlsVerify:             true,
	storageDriver:         "vfs",
	format:                "oci",
	dockerfile:            "./Dockerfile",
	contextDir:            "docker",
	nexusURL:              os.Getenv("NEXUS_URL"),
	nexusUsername:         os.Getenv("NEXUS_USERNAME"),
	nexusPassword:         os.Getenv("NEXUS_PASSWORD"),
	buildahBuildExtraArgs: "",
	buildahPushExtraArgs:  "",
	trivySBOMExtraArgs:    "",
	aquasecGate:           false,
	debug:                 (os.Getenv("DEBUG") == "true"),
	sbomFormat:            "spdx",
}

func main() {
	opts := options{}
	flag.StringVar(&opts.checkoutDir, "checkout-dir", defaultOptions.checkoutDir, "Checkout dir")
	flag.StringVar(&opts.bitbucketAccessToken, "bitbucket-access-token", defaultOptions.bitbucketAccessToken, "bitbucket-access-token")
	flag.StringVar(&opts.bitbucketURL, "bitbucket-url", defaultOptions.bitbucketURL, "bitbucket-url")
	flag.StringVar(&opts.aquaUsername, "aqua-username", defaultOptions.aquaUsername, "aqua-username")
	flag.StringVar(&opts.aquaPassword, "aqua-password", defaultOptions.aquaPassword, "aqua-password")
	flag.StringVar(&opts.aquaURL, "aqua-url", defaultOptions.aquaURL, "aqua-url")
	flag.StringVar(&opts.aquaRegistry, "aqua-registry", defaultOptions.aquaRegistry, "aqua-registry")
	flag.StringVar(&opts.imageStream, "image-stream", defaultOptions.imageStream, "Image stream")
	flag.StringVar(&opts.registry, "registry", defaultOptions.registry, "Registry")
	flag.StringVar(&opts.certDir, "cert-dir", defaultOptions.certDir, "Use certificates at the specified path to access the registry")
	flag.StringVar(&opts.imageNamespace, "image-namespace", defaultOptions.imageNamespace, "image namespace")
	flag.BoolVar(&opts.tlsVerify, "tls-verify", defaultOptions.tlsVerify, "TLS verify")
	flag.StringVar(&opts.storageDriver, "storage-driver", defaultOptions.storageDriver, "storage driver")
	flag.StringVar(&opts.format, "format", defaultOptions.format, "format of the built container, oci or docker")
	flag.StringVar(&opts.dockerfile, "dockerfile", defaultOptions.dockerfile, "dockerfile")
	flag.StringVar(&opts.contextDir, "context-dir", defaultOptions.contextDir, "contextDir")
	flag.StringVar(&opts.nexusURL, "nexus-url", defaultOptions.nexusURL, "Nexus URL")
	flag.StringVar(&opts.nexusUsername, "nexus-username", defaultOptions.nexusUsername, "Nexus username")
	flag.StringVar(&opts.nexusPassword, "nexus-password", defaultOptions.nexusPassword, "Nexus password")
	flag.StringVar(&opts.buildahBuildExtraArgs, "buildah-build-extra-args", defaultOptions.buildahBuildExtraArgs, "extra parameters passed for the build command when building images")
	flag.StringVar(&opts.buildahPushExtraArgs, "buildah-push-extra-args", defaultOptions.buildahPushExtraArgs, "extra parameters passed for the push command when pushing images")
	flag.StringVar(&opts.trivySBOMExtraArgs, "trivy-sbom-extra-args", defaultOptions.trivySBOMExtraArgs, "extra parameters passed for the trivy command to generate an SBOM")
	flag.StringVar(&opts.sbomFormat, "sbom-format", defaultOptions.sbomFormat, "SBOM format")
	flag.BoolVar(&opts.aquasecGate, "aqua-gate", defaultOptions.aquasecGate, "whether the Aqua security scan needs to pass for the task to succeed")
	flag.BoolVar(&opts.debug, "debug", defaultOptions.debug, "debug mode")
	flag.Parse()

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	err := (&packageImage{logger: logger, opts: opts}).runSteps(
		setupContext(),
		setImageName(),
		skipIfImageDigestExists(),
		buildImageAndGenerateTar(),
		generateSBOM(),
		pushImage(),
		scanImageWithAqua(),
		storeArtifact(),
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

// copyAquaReportsToArtifacts copies the Aqua scan reports to the artifacts directory.
func copyAquaReportsToArtifacts(htmlReportFile, jsonReportFile string) error {
	if _, err := os.Stat(htmlReportFile); err == nil {
		err := pipelinectxt.CopyArtifact(htmlReportFile, pipelinectxt.AquaScansPath)
		if err != nil {
			return fmt.Errorf("copying HTML report to artifacts failed: %w", err)
		}
	}
	if _, err := os.Stat(jsonReportFile); err == nil {
		err := pipelinectxt.CopyArtifact(jsonReportFile, pipelinectxt.AquaScansPath)
		if err != nil {
			return fmt.Errorf("copying JSON report to artifacts failed: %w", err)
		}
	}
	return nil
}

// createBitbucketInsightReport attaches a code insight report to the Git commit
// being built in Bitbucket. The code insight report points to the Aqua security scan.
func createBitbucketInsightReport(opts options, aquaScanUrl string, success bool, ctxt *pipelinectxt.ODSContext) error {
	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	}
	bitbucketClient, err := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: opts.bitbucketAccessToken,
		BaseURL:  opts.bitbucketURL,
		Logger:   logger,
	})
	if err != nil {
		return fmt.Errorf("bitbucket client: %w", err)
	}
	reportKey := "org.opendevstack.aquasec"
	scanResult := bitbucket.InsightReportFail
	if success {
		scanResult = bitbucket.InsightReportPass
	}
	_, err = bitbucketClient.InsightReportCreate(
		ctxt.Project,
		ctxt.Repository,
		ctxt.GitCommitSHA,
		reportKey,
		bitbucket.InsightReportCreatePayload{
			Title:       "Aqua Security",
			Reporter:    "OpenDevStack",
			CreatedDate: time.Now().Unix(),
			Details:     "Please visit the following link to review the Aqua Security scan report:",
			Result:      scanResult,
			Data: []bitbucket.InsightReportData{
				{
					Title: "Report",
					Type:  "LINK",
					Value: map[string]string{
						"linktext": "Result in Aqua",
						"href":     aquaScanUrl,
					},
				},
			},
		},
	)
	return err
}

// getImageDigestFromFile reads the image digest from the file written to by buildah.
func getImageDigestFromFile(workingDir string) (string, error) {
	content, err := os.ReadFile(filepath.Join(workingDir, "image-digest"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

// writeImageDigestToResults writes the image digest into the Tekton results file.
func writeImageDigestToResults(imageDigest string) error {
	err := os.MkdirAll("/tekton/results", 0644)
	if err != nil {
		return err
	}
	return os.WriteFile("/tekton/results/image-digest", []byte(imageDigest), 0644)
}
