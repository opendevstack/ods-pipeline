package main

import (
	"flag"
	"os"

	"github.com/opendevstack/pipeline/internal/image"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"golang.org/x/exp/slog"
)

type options struct {
	checkoutDir          string
	imageStream          string
	imageNamespace       string
	bitbucketAccessToken string
	bitbucketURL         string
	aquaUsername         string
	aquaPassword         string
	aquaURL              string
	aquaRegistry         string
	aquasecGate          bool
	debug                bool
}

type aquaScan struct {
	opts    options
	ctxt    *pipelinectxt.ODSContext
	imageId image.Identity
}

var defaultOptions = options{
	checkoutDir:          ".",
	imageStream:          "",
	imageNamespace:       "",
	bitbucketAccessToken: os.Getenv("BITBUCKET_ACCESS_TOKEN"),
	bitbucketURL:         os.Getenv("BITBUCKET_URL"),
	aquaUsername:         os.Getenv("AQUA_USERNAME"),
	aquaPassword:         os.Getenv("AQUA_PASSWORD"),
	aquaURL:              os.Getenv("AQUA_URL"),
	aquaRegistry:         os.Getenv("AQUA_REGISTRY"),
	aquasecGate:          false,
	debug:                (os.Getenv("DEBUG") == "true"),
}

func main() {
	opts := options{}
	flag.StringVar(&opts.checkoutDir, "checkout-dir", defaultOptions.checkoutDir, "Checkout dir")
	flag.StringVar(&opts.imageStream, "image-stream", defaultOptions.imageStream, "Image stream")
	flag.StringVar(&opts.imageNamespace, "image-namespace", defaultOptions.imageNamespace, "image namespace")
	flag.StringVar(&opts.bitbucketAccessToken, "bitbucket-access-token", defaultOptions.bitbucketAccessToken, "bitbucket-access-token")
	flag.StringVar(&opts.bitbucketURL, "bitbucket-url", defaultOptions.bitbucketURL, "bitbucket-url")
	flag.StringVar(&opts.aquaUsername, "aqua-username", defaultOptions.aquaUsername, "aqua-username")
	flag.StringVar(&opts.aquaPassword, "aqua-password", defaultOptions.aquaPassword, "aqua-password")
	flag.StringVar(&opts.aquaURL, "aqua-url", defaultOptions.aquaURL, "aqua-url")
	flag.StringVar(&opts.aquaRegistry, "aqua-registry", defaultOptions.aquaRegistry, "aqua-registry")
	flag.BoolVar(&opts.aquasecGate, "aqua-gate", defaultOptions.aquasecGate, "whether the Aqua security scan needs to pass for the task to succeed")
	flag.BoolVar(&opts.debug, "debug", defaultOptions.debug, "debug mode")
	flag.Parse()

	logLevel := slog.LevelInfo
	if opts.debug {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.HandlerOptions{Level: logLevel}.NewTextHandler(os.Stderr)))

	err := (&aquaScan{opts: opts}).runSteps(
		setupContext(),
		setImageId(),
		skipIfScanArtifactsExist(),
		scanImagesWithAqua(),
	)
	if err != nil {
		slog.Error("step failed", err)
		os.Exit(1)
	}
}
