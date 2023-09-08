package tektontaskrun

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
)

const (
	DefaultServiceAccountName = "pipeline"
	KinDMountHostPath         = "/tmp/ods-pipeline/kind-mount"
	KinDMountContainerPath    = "/files"
	KinDRegistry              = "localhost:5000"
	KinDName                  = "ods-pipeline"
)

var recreateClusterFlag = flag.Bool("ods-recreate-cluster", false, "Whether to remove and recreate the KinD cluster named 'ods-pipeline'")
var registryPortFlag = flag.String("ods-cluster-registry-port", "5000", "Port of cluster registry")
var outsideKindFlag = flag.Bool("ods-outside-kind", false, "Whether to continue if the Kube context is not set to the KinD cluster")
var reuseImagesFlag = flag.Bool("ods-reuse-images", false, "Whether to reuse existing images instead of building again")
var debugFlag = flag.Bool("ods-debug", false, "Turn on debug mode for scripts etc.")

// ClusterOpt allows to further configure the KinD cluster after its creation.
type ClusterOpt func(c *ClusterConfig) error

// ClusterConfig represents key configuration of the KinD cluster.
type ClusterConfig struct {
	StorageSourceDir  string
	StorageCapacity   string
	StorageClassName  string
	Registry          string
	DefaultRepository string
}

// ImageBuildConfig represents the config used to build a container image.
type ImageBuildConfig struct {
	Dockerfile string
	Tag        string
	ContextDir string
}

// Process validates the configuration and defaults the image tag if unset
// using the defaultImageRepository and Dockerfile values.
func (ibc *ImageBuildConfig) Process(defaultImageRepository string) error {
	if ibc.Dockerfile == "" || ibc.ContextDir == "" {
		return errors.New("both Dockerfile and ContextDir must be set")
	}
	if ibc.Tag == "" {
		imageName := strings.TrimPrefix(path.Base(ibc.Dockerfile), "Dockerfile.")
		ibc.Tag = fmt.Sprintf("%s/%s:latest", defaultImageRepository, imageName)
	}
	return nil
}

// NewClusterConfig creates a new ClusterConfig instance.
func NewClusterConfig() *ClusterConfig {
	return &ClusterConfig{
		StorageClassName:  "standard", // if using KinD, set it to "standard"
		StorageCapacity:   "1Gi",
		StorageSourceDir:  KinDMountContainerPath,
		Registry:          KinDRegistry,
		DefaultRepository: "ods-pipeline",
	}
}

// DefaultImageRepository returns the registry + default repository
// combination.
func (c *ClusterConfig) DefaultImageRepository() string {
	return c.Registry + "/" + c.DefaultRepository
}

// DefaultTaskTemplateData returns a map with default values which can be used
// in task templates.
func (c *ClusterConfig) DefaultTaskTemplateData() map[string]string {
	return map[string]string{
		"ImageRepository": c.DefaultImageRepository(),
		"Version":         "latest",
	}
}

// StartKinDCluster starts a KinD cluster with Tekton installed.
// Afterwards, any given ClusterOpt is applied.
func StartKinDCluster(opts ...ClusterOpt) (*ClusterConfig, error) {
	flag.Parse()
	if err := checkCluster(*outsideKindFlag); err != nil {
		return nil, fmt.Errorf("check kubectl context: %s", err)
	}
	if err := createKinDCluster(*debugFlag); err != nil {
		return nil, fmt.Errorf("create KinD cluster: %s", err)
	}
	if err := installTektonPipelines(*debugFlag); err != nil {
		return nil, fmt.Errorf("install Tekton: %s", err)
	}

	c := NewClusterConfig()
	for _, o := range opts {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// LoadImage builds a container image using the docker CLI based on the given
// ImageBuildConfig.
//
// The ImageBuildConfig must set at least Dockerfile and ContextDir option.
// If Tag is unset, it is inferred from the default registry and the Dockerfile
// name. For example, given a Dockerfile of "Dockerfile.foobar", the tag is
// defaulted to localhost:5000/ods-pipeline/foobar.
// Passing the flag -ods-reuse-images to the tests will skip image rebuilding.
func LoadImage(ibc ImageBuildConfig) ClusterOpt {
	flag.Parse()
	return func(c *ClusterConfig) error {
		buildImage := true
		err := ibc.Process(c.DefaultImageRepository())
		if err != nil {
			return fmt.Errorf("processing image build config: %s", err)
		}
		if *reuseImagesFlag {
			cmd := exec.Command("docker", "images", "-q", ibc.Tag)
			b, err := cmd.Output()
			if err != nil {
				return err
			}
			imageID := strings.TrimSpace(string(b))
			if imageID != "" {
				log.Printf("Reusing image ID %s for tag %s ...\n", imageID, ibc.Tag)
				buildImage = false
			}
		}
		if buildImage {
			log.Printf("Building image %s from %s ...\n", ibc.Tag, ibc.Dockerfile)
			if !path.IsAbs(ibc.Dockerfile) {
				ibc.Dockerfile = filepath.Join(ibc.ContextDir, ibc.Dockerfile)
			}
			args := []string{
				"build",
				"-f", ibc.Dockerfile,
				"-t", ibc.Tag,
				ibc.ContextDir,
			}
			if err := command.Run("docker", args, []string{}, os.Stdout, os.Stderr); err != nil {
				return err
			}
		}
		return command.Run("docker", []string{"push", ibc.Tag}, []string{}, os.Stdout, os.Stderr)
	}
}

func checkCluster(outsideKindAllowed bool) error {
	if !outsideKindAllowed {
		cmd := exec.Command("kubectl", "config", "current-context")
		b, err := cmd.Output()
		if err != nil || len(b) == 0 {
			log.Println("did not detect existing kubectl context")
			return nil
		}
		gotContext := strings.TrimSpace(string(b))
		wantCluster := "ods-pipeline"
		if gotContext != "kind-"+wantCluster {
			return fmt.Errorf("not running tests outside KinD cluster ('%s') without -ods-outside-kind! Current context: %s", wantCluster, gotContext)
		}
	}
	return nil
}

func createKinDCluster(debug bool) error {
	args := []string{
		projectpath.RootedPath("scripts/kind-with-registry.sh"),
		"--registry-port=" + *registryPortFlag,
	}
	if *recreateClusterFlag {
		args = append(args, "--recreate")
	}
	if debug {
		args = append(args, "--verbose")
	}
	return command.Run("bash", args, []string{}, os.Stdout, os.Stderr)
}

func installTektonPipelines(debug bool) error {
	args := []string{
		projectpath.RootedPath("scripts/install-tekton-pipelines.sh"),
	}
	if debug {
		args = append(args, "--verbose")
	}
	return command.Run("sh", args, []string{}, os.Stdout, os.Stderr)
}
