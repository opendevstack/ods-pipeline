package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/image"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

const (
	kubernetesServiceaccountDir = "/var/run/secrets/kubernetes.io/serviceaccount"
)

type options struct {
	checkoutDir           string
	imageStream           string
	extraTags             string
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
	debug                 bool
}

type packageImage struct {
	logger          logging.LeveledLoggerInterface
	opts            options
	parsedExtraTags []string
	ctxt            *pipelinectxt.ODSContext
	imageId         image.Identity
	imageDigest     string
}

func (p *packageImage) imageName() string {
	return p.imageId.StreamSha()
}

func (p *packageImage) imageNameNoSha() string {
	return p.imageId.ImageStream
}

func (p *packageImage) imageRef() string {
	return p.imageId.ImageRefWithSha(p.opts.registry)
}

func (p *packageImage) artifactImage() artifact.Image {
	return p.imageId.ArtifactImage(p.opts.registry, p.imageDigest)
}

func (p *packageImage) artifactImageForTag(tag string) artifact.Image {
	imageExtraTag := p.imageId.Tag(tag)
	return imageExtraTag.ArtifactImage(p.opts.registry, p.imageDigest)
}

var defaultOptions = options{
	checkoutDir:           ".",
	imageStream:           "",
	extraTags:             "",
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
	debug:                 (os.Getenv("DEBUG") == "true"),
}

func main() {
	opts := options{}
	flag.StringVar(&opts.checkoutDir, "checkout-dir", defaultOptions.checkoutDir, "Checkout dir")
	flag.StringVar(&opts.imageStream, "image-stream", defaultOptions.imageStream, "Image stream")
	flag.StringVar(&opts.extraTags, "extra-tags", defaultOptions.extraTags, "Extra tags")
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
	flag.BoolVar(&opts.debug, "debug", defaultOptions.debug, "debug mode")
	flag.Parse()
	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}
	p := packageImage{logger: logger, opts: opts}
	err := (&p).runSteps(
		setExtraTags(),
		setupContext(),
		setImageId(),
		skipIfImageArtifactExists(),
		buildImageAndGenerateTar(),
		generateSBOM(),
		pushImage(),
		storeArtifact(),
	)
	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
	// If skipIfImageArtifactExists skips the remaining runSteps, extra-tags
	// still should be processed if their related artifact has not been set.
	err = (&p).runSteps(processExtraTags())
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

// imageArtifactExists checks if image artifact JSON file exists in its artifacts path
func imageArtifactExists(p *packageImage) error {
	imageArtifactsDir := filepath.Join(p.opts.checkoutDir, pipelinectxt.ImageDigestsPath)
	imageArtifactFilename := fmt.Sprintf("%s.json", p.ctxt.Component)
	_, err := os.Stat(filepath.Join(imageArtifactsDir, imageArtifactFilename))
	return err
}
