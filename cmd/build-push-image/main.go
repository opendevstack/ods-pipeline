package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

const (
	aquasecBin = "aquasec"
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
	imageNamespace        string
	tlsVerify             bool
	storageDriver         string
	format                string
	dockerfile            string
	contextDir            string
	buildahBuildExtraArgs string
	buildahPushExtraArgs  string
	aquasecGate           bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.bitbucketAccessToken, "bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	flag.StringVar(&opts.bitbucketURL, "bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	flag.StringVar(&opts.aquaUsername, "aqua-username", os.Getenv("AQUA_USERNAME"), "aqua-username")
	flag.StringVar(&opts.aquaPassword, "aqua-password", os.Getenv("AQUA_PASSWORD"), "aqua-password")
	flag.StringVar(&opts.aquaURL, "aqua-url", os.Getenv("AQUA_URL"), "aqua-url")
	flag.StringVar(&opts.aquaRegistry, "aqua-registry", os.Getenv("AQUA_REGISTRY"), "aqua-registry")
	flag.StringVar(&opts.imageStream, "image-stream", "", "Image stream")
	flag.StringVar(&opts.registry, "registry", "image-registry.openshift-image-registry.svc:5000", "Registry")
	flag.StringVar(&opts.imageNamespace, "image-namespace", "", "image namespace")
	flag.BoolVar(&opts.tlsVerify, "tls-verify", true, "TLS verify")
	flag.StringVar(&opts.storageDriver, "storage-driver", "vfs", "storage driver")
	flag.StringVar(&opts.format, "format", "oci", "format of the built container, oci or docker")
	flag.StringVar(&opts.dockerfile, "dockerfile", "./Dockerfile", "dockerfile")
	flag.StringVar(&opts.contextDir, "context-dir", "docker", "contextDir")
	flag.StringVar(&opts.buildahBuildExtraArgs, "buildah-build-extra-args", "docker", "extra parameters passed for the build command when building images")
	flag.StringVar(&opts.buildahPushExtraArgs, "buildah-push-extra-args", "docker", "extra parameters passed for the push command when pushing images")
	flag.BoolVar(&opts.aquasecGate, "aqua-gate", false, "whether the Aqua security scan needs to pass for the task to succeed")
	flag.Parse()

	workingDir := "."
	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(workingDir)
	if err != nil {
		log.Fatal(err)
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
	imageDigest, err := getImageDigestFromRegistry(imageRef, opts.tlsVerify)
	if err == nil {
		fmt.Println("Image exists already.")
	} else {

		fmt.Printf("Building image %s ...\n", imageName)
		stdout, stderr, err := buildahBuild(opts, imageRef)
		if err != nil {
			log.Println(string(stderr))
			log.Fatal(err)
		}
		fmt.Println(string(stdout))

		fmt.Printf("Pushing image %s ...\n", imageRef)
		stdout, stderr, err = buildahPush(opts, workingDir, imageRef)
		if err != nil {
			log.Println(string(stderr))
			log.Fatal(err)
		}
		fmt.Println(string(stdout))

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
			scanSuccessful := false

			stdout, stderr, err = aquaScan(opts, aquaImage, htmlReportFile, jsonReportFile)
			if err != nil {
				log.Println(string(stderr))
				log.Println(err)
				if opts.aquasecGate {
					log.Fatalln("Stopping build as successful Aqua scan is required")
				}
			} else {
				scanSuccessful = true
			}
			fmt.Println(string(stdout))

			aquaScanUrl := fmt.Sprintf(
				"%s/#/images/%s/%s/vulns",
				opts.aquaURL, url.QueryEscape(opts.aquaRegistry), url.QueryEscape(aquaImage),
			)
			fmt.Printf("Aqua vulnerability report is at %s ...\n", aquaScanUrl)

			err = copyAquaReportsToArtifacts(htmlReportFile, jsonReportFile)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Creating Bitbucket code insight report ...")
			err = createBitbucketInsightReport(opts, aquaScanUrl, scanSuccessful, ctxt)
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

// buildahBuild builds a local image using the Dockerfile and context directory
// given in opts, tagging the resulting image with given tag.
func buildahBuild(opts options, tag string) ([]byte, []byte, error) {
	return command.Run("buildah", []string{
		fmt.Sprintf("--storage-driver=%s", opts.storageDriver),
		"bud",
		fmt.Sprintf("--format=%s", opts.format),
		fmt.Sprintf("--tls-verify=%v", opts.tlsVerify),
		"--no-cache",
		fmt.Sprintf("--file=%s", opts.dockerfile),
		fmt.Sprintf("--tag=%s", tag),
		opts.contextDir,
	})
}

// buildahPush pushes a local image to the given imageRef.
func buildahPush(opts options, workingDir, imageRef string) ([]byte, []byte, error) {
	return command.Run("buildah", []string{
		fmt.Sprintf("--storage-driver=%s", opts.storageDriver),
		"push",
		fmt.Sprintf("--tls-verify=%v", opts.tlsVerify),
		fmt.Sprintf("--digestfile=%s", filepath.Join(workingDir, "image-digest")),
		imageRef, fmt.Sprintf("docker://%s", imageRef),
	})
}

// buildahPush pushes a local image to the given imageName.
// func buildahOCIArchive(opts options, workingDir, imageName string) ([]byte, []byte, error) {
// 	return command.Run("buildah", []string{
// 		fmt.Sprintf("--storage-driver=%s", opts.storageDriver),
// 		"push",
// 		imageName, fmt.Sprintf("oci-archive:/%s:%s", filepath.Join(workingDir, "image.tar"), imageName),
// 	})
// }

// aquaScan runs an Aqua scan on given image.
func aquaScan(opts options, image, htmlReportFile, jsonReportFile string) ([]byte, []byte, error) {
	return command.Run(aquasecBin, []string{
		"scan",
		"--dockerless", "--register", "--text",
		fmt.Sprintf("--htmlfile=%s", htmlReportFile),
		fmt.Sprintf("--jsonfile=%s", jsonReportFile),
		"-w", "/tmp",
		fmt.Sprintf("--user=%s", opts.aquaUsername),
		fmt.Sprintf("--password=%s", opts.aquaPassword),
		fmt.Sprintf("--host=%s", opts.aquaURL),
		image,
		fmt.Sprintf("--registry=%s", opts.aquaRegistry),
	})
}

// aquaImport runs an Aqua scan on given imageRef.
// func aquaImport(opts options, workingDir, jsonReportFile string) ([]byte, []byte, error) {
// 	return command.Run(aquasecBin, []string{
// 		"import",
// 		fmt.Sprintf("--user=%s", opts.aquaUsername),
// 		fmt.Sprintf("--password=%s", opts.aquaPassword),
// 		fmt.Sprintf("--host=%s", opts.aquaURL),
// 		jsonReportFile,
// 		"-w", "/tmp",
// 	})
// }

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
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: opts.bitbucketAccessToken,
		BaseURL:  opts.bitbucketURL,
	})
	reportKey := "org.opendevstack.aquasec"
	scanResult := bitbucket.InsightReportFail
	if success {
		scanResult = bitbucket.InsightReportPass
	}
	_, err := bitbucketClient.InsightReportCreate(
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
					Title: "Link",
					Type:  "LINK",
					Value: map[string]string{
						"linktext": aquaScanUrl,
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
func getImageDigestFromRegistry(imageRef string, tlsVerify bool) (string, error) {
	stdout, _, err := command.Run("skopeo", []string{
		"inspect",
		fmt.Sprintf("--format=%s", "{{.Digest}}"),
		fmt.Sprintf("--tls-verify=%v", tlsVerify),
		"docker://" + imageRef,
	})
	return string(stdout), err
}

// getImageDigestFromFile reads the image digest from the file written to by buildah.
func getImageDigestFromFile(workingDir string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(workingDir, "image-digest"))
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
	return ioutil.WriteFile("/tekton/results/image-digest", []byte(imageDigest), 0644)
}

// aquasecInstalled checks whether the Aqua binary is in the $PATH.
func aquasecInstalled() bool {
	_, err := exec.LookPath(aquasecBin)
	return err == nil
}
