package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

// Stage is a stage identifier. There are three stages: DEV, QA and PROD.
// Note that stages are not environments but rather describe what kind of
// constraints apply to a certain environment. For example, environments
// of stage PROD are protected from deployment without a prior QA deployment.
type Stage string

const (
	DevStage      = "dev"
	QAStage       = "qa"
	ProdStage     = "prod"
	DefaultBranch = "refs/heads/master"
	ODSYMLFile    = "ods.yml"
	ODSYAMLFile   = "ods.yaml"
)

var ODSFileCandidates = []string{ODSYAMLFile, ODSYMLFile}

// ODS represents the ODS pipeline configuration for one repository.
type ODS struct {
	// Repositories specifies the subrepositores, making the current repository
	// an "umbrella" repository.
	Repositories []Repository `json:"repositories,omitempty"`
	// Environments allows you to specify target environments to deploy to.
	Environments []Environment `json:"environments"`
	// BranchToEnvironmentMapping configures which branch should be deployed to which environment.
	BranchToEnvironmentMapping []BranchToEnvironmentMapping `json:"branchToEnvironmentMapping,omitempty"`
	// Pipeline allows to define the Tekton pipeline tasks.
	Pipeline Pipeline `json:"pipeline,omitempty"`
	// Version is the application version and must follow SemVer.
	Version string `json:"version,omitempty"`
}

// Repository represents a Git repository.
type Repository struct {
	// Name of the Git repository (without host/organisation and trailing .git)
	// Example: "foobar"
	Name string `json:"name"`
	// URL of Git repository (optional). If not given, the repository given by
	// Name is assumed to be under the same organisation than the repository
	// hosting the ods.y(a)ml file.
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
	// API server of the target cluster, including scheme.
	APIServer string `json:"apiServer"`
	// Name of the Secret resource holding the API user credentials.
	APICredentialsSecret string `json:"apiCredentialsSecret"`
	// Hostname of the target registry. If not given, the registy of the source
	// image is used.
	RegistryHost string `json:"registryHost"`
	// Whether to verify TLS for the target registry.
	RegistryTLSVerify *bool `json:"registryVerify,omitempty"`
	// Target K8s namespace (OpenShift project) on the target cluster to deploy into.
	Namespace string `json:"namespace"`
	// Additional configuration of the target. This may be used by tasks outside
	// the ODS catalog.
	Config map[string]interface{} `json:"config,omitempty"`
	// APIToken holds the token of the environment, if any.
	// The value is retrieved from the "token" field in the secret referenced by APICredentialsSecret.
	// Cannot be set from JSON.
	APIToken string `json:"-"`
}

// Pipeline represents a Tekton pipeline.
type Pipeline struct {
	Tasks   []tekton.PipelineTask `json:"tasks,omitempty"`
	Finally []tekton.PipelineTask `json:"finally,omitempty"`
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
	if len(e.Name) == 0 {
		return errors.New("name of environment must not be blank")
	}
	switch e.Stage {
	case DevStage, QAStage, ProdStage:
		return nil
	default:
		return fmt.Errorf("invalid stage value '%s' for environment %s", e.Stage, e.Name)
	}
}

func (o *ODS) Environment(environment string) (*Environment, error) {
	var envs []string
	for _, e := range o.Environments {
		if e.Name == environment {
			return &e, nil
		}
		envs = append(envs, e.Name)
	}

	return nil, fmt.Errorf("no environment matched '%s', have: %s", environment, strings.Join(envs, ", "))
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
	for _, c := range ODSFileCandidates {
		candidate := filepath.Join(dir, c)
		if _, err := os.Stat(candidate); err == nil {
			return ReadFromFile(candidate)
		}
	}
	return nil, fmt.Errorf("no matching file in '%s', looked for: %s", dir, strings.Join(ODSFileCandidates, ", "))
}
