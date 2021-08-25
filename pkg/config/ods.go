package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

type Stage string

const (
	Dev           Stage  = "dev"
	QA                   = "qa"
	Prod                 = "prod"
	DefaultBranch string = "refs/heads/master"
)

type ODS struct {
	Repositories               []Repository                 `json:"repositories"`
	Environments               []Environment                `json:"environments"`
	BranchToEnvironmentMapping []BranchToEnvironmentMapping `json:"branchToEnvironmentMapping"`
	Pipeline                   Pipeline                     `json:"pipeline"`
}

// Repository represents a Git repository.
type Repository struct {
	// Name of the Git repository (without host/organisation and trailing .git)
	// Example: "foobar"
	Name string `json:"name"`
	// URL of Git repository (optional). If not given, the repository given by
	// Name is assumed to be under the same organisation than the repository
	// hosting the ods.yml file.
	// Example: "https://acme.org/foo/bar.git"
	URL string `json:"url"`
	// Branch of Git repository (optional). If none is given, this defaults to
	// the "master" branch.
	// Example: "develop"
	Branch string `json:"branch"`
}

type BranchToEnvironmentMapping struct {
	// Name of Git branch. May also be a prefix like "release/*"
	Branch string `json:"branch"`
	// Environment of the environment.
	Environment string `json:"environment"`
}

type Environment struct {
	// Name of the environment to deploy to. This is an arbitary name.
	Name string `json:"name"`
	// Kind of the environment to deploy to. One of "dev", "qa", "prod".
	Stage Stage `json:"stage"`
	// API URL of the target cluster.
	URL string `json:"url"`
	// Hostname of the target registry. If not given, the registy of the source
	// image is used.
	RegistryHost string `json:"registryHost"`
	// Whether to verify TLS for the target registry.
	RegistryTLSVerify *bool `json:"registryVerify"`
	// Target K8s namespace (OpenShift project) on the target cluster to deploy into.
	Namespace string `json:"namespace"`
	// Name of the Secret resource holding the API user credentials.
	SecretRef string `json:"secretRef"`
	// Additional configuration of the target. This may be used by tasks outside
	// the ODS catalog.
	Config map[string]interface{} `json:"config"`
}

func (o *ODS) Validate() error {
	for _, e := range o.Environments {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (e Environment) Validate() error {
	switch e.Stage {
	case Dev, QA, Prod:
		return nil
	default:
		return fmt.Errorf("invalid stage value '%s' for environment %s", e.Stage, e.Name)
	}
}

// Pipeline represents a Tekton pipeline.
type Pipeline struct {
	Tasks   []tekton.PipelineTask `json:"tasks"`
	Finally []tekton.PipelineTask `json:"finally"`
}

// Read reads an ods config from given byte slice or errors.
func Read(body []byte) (*ODS, error) {
	var odsConfig *ODS
	err := yaml.UnmarshalStrict(body, &odsConfig, func(dec *json.Decoder) *json.Decoder {
		dec.DisallowUnknownFields()
		return dec
	})
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	if err = odsConfig.Validate(); err != nil {
		return nil, err
	}
	return odsConfig, nil
}

// ReadFromFile reads an ods config from given filename or errors.
func ReadFromFile(filename string) (*ODS, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}
	return Read(body)
}

// ReadFromDir reads an ods config file from given dir or errors.
func ReadFromDir(dir string) (*ODS, error) {
	candidates := []string{
		filepath.Join(dir, "ods.yml"),
		filepath.Join(dir, "ods.yaml"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return ReadFromFile(c)
		}
	}
	return nil, fmt.Errorf("no matching file in '%s', looked for: %s", dir, strings.Join(candidates, ", "))
}
