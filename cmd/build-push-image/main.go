package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

const (
	kubernetesServiceaccountDir = "/var/run/secrets/kubernetes.io/serviceaccount"
)

type options struct {
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
	aquasecGate           bool
	debug                 bool
}

var defaultOptions = options{
	bitbucketAccessToken:  os.Getenv("BITBUCKET_ACCESS_TOKEN"),
	bitbucketURL:          os.Getenv("BITBUCKET_URL"),
	aquaUsername:          os.Getenv("AQUA_USERNAME"),
	aquaPassword:          os.Getenv("AQUA_PASSWORD"),
	aquaURL:               os.Getenv("AQUA_URL"),
	aquaRegistry:          os.Getenv("AQUA_REGISTRY"),
	imageStream:           "",
	registry:              "image-registry.openshift-image-registry.svc:5000",
	certDir:               "/etc/containers/certs.d",
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
	aquasecGate:           false,
	debug:                 (os.Getenv("DEBUG") == "true"),
}

func main() {
	opts := options{}
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
	flag.BoolVar(&opts.aquasecGate, "aqua-gate", defaultOptions.aquasecGate, "whether the Aqua security scan needs to pass for the task to succeed")
	flag.BoolVar(&opts.debug, "debug", defaultOptions.debug, "debug mode")
	flag.Parse()

	workingDir := "."
	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(workingDir)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(kubernetesServiceaccountDir); err == nil {
		opts.certDir = kubernetesServiceaccountDir
	}
	if opts.debug {
		directory.ListFiles(opts.certDir)
	}

	// TLS verification of the KinD registry is not possible at the moment as
	// requests error out with "server gave HTTP response to HTTPS client".
	if strings.HasPrefix(opts.registry, "kind-registry.kind") {
		opts.tlsVerify = false
	}

	imageNamespace := opts.imageNamespace
	if len(imageNamespace) == 0 {
		imageNamespace = ctxt.Namespace
	}
	imageStream := opts.imageStream
	if len(imageStream) == 0 {
		imageStream = ctxt.Component
	}
	imageName := fmt.Sprintf("%s:%s", imageStream, ctxt.GitCommitSHA)
	imageRef := fmt.Sprintf(
		"%s/%s/%s",
		opts.registry, imageNamespace, imageName,
	)

	fmt.Printf("Checking if image %s exists already ...\n", imageName)
	imageDigest, err := getImageDigestFromRegistry(imageRef, opts)
	if err == nil {
		fmt.Println("Image exists already.")
	} else {

		fmt.Printf("Building image %s ...\n", imageName)
		err := buildahBuild(opts, imageRef, os.Stdout, os.Stderr)
		if err != nil {
			log.Fatal("buildah bud: ", err)
		}

		fmt.Printf("Pushing image %s ...\n", imageRef)
		err = buildahPush(opts, workingDir, imageRef, os.Stdout, os.Stderr)
		if err != nil {
			log.Fatal("buildah push: ", err)
		}

		d, err := getImageDigestFromFile(workingDir)
		if err != nil {
			log.Fatal(err)
		}
		imageDigest = d

		if aquasecInstalled() {
			fmt.Println("Scanning image with Aqua scanner ...")
			aquaImage := fmt.Sprintf("%s/%s", imageNamespace, imageName)
			htmlReportFile := filepath.Join(workingDir, "report.html")
			jsonReportFile := filepath.Join(workingDir, "report.json")
			scanArgs := aquaAssembleScanArgs(opts, aquaImage, htmlReportFile, jsonReportFile)
			scanSuccessful, err := aquaScan(aquasecBin, scanArgs, os.Stdout, os.Stderr)
			if err != nil {
				log.Fatal("aqua scan: ", err)
			}

			if !scanSuccessful && opts.aquasecGate {
				log.Fatalln("Stopping build as successful Aqua scan is required")
			}

			asu, err := aquaScanURL(opts, aquaImage)
			if err != nil {
				log.Fatal("aqua scan URL:", err)
			}
			fmt.Printf("Aqua vulnerability report is at %s ...\n", asu)

			err = copyAquaReportsToArtifacts(htmlReportFile, jsonReportFile)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Creating Bitbucket code insight report ...")
			err = createBitbucketInsightReport(opts, asu, scanSuccessful, ctxt)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Aqua is not configured, image will not be scanned for vulnerabilities.")
		}
	}

	err = writeImageDigestToResults(imageDigest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Writing image artifact ...")
	ia := artifact.Image{
		Image:      imageRef,
		Registry:   opts.registry,
		Repository: imageNamespace,
		Name:       imageStream,
		Tag:        ctxt.GitCommitSHA,
		Digest:     imageDigest,
	}
	imageArtifactFilename := fmt.Sprintf("%s.json", imageStream)
	err = pipelinectxt.WriteJsonArtifact(ia, pipelinectxt.ImageDigestsPath, imageArtifactFilename)
	if err != nil {
		log.Fatal(err)
	}
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

// getImageDigestFromRegistry returns a SHA256 image digest if the specified
// imageRef exists. Example return value:
// "sha256:3b6de1c737065e9973ddb7cc60b769b866b7649ff6f2de3816934dda832de294"
func getImageDigestFromRegistry(imageRef string, opts options) (string, error) {
	args := []string{
		"inspect",
		fmt.Sprintf("--format=%s", "{{.Digest}}"),
		fmt.Sprintf("--tls-verify=%v", opts.tlsVerify),
		fmt.Sprintf("--cert-dir=%s", opts.certDir),
	}
	if opts.debug {
		args = append(args, "--debug")
	}
	stdout, _, err := command.Run("skopeo", append(args, "docker://"+imageRef))
	return string(stdout), err
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
