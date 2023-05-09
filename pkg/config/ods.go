package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/pod"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
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
	// Pipeline allows to define the Tekton pipeline tasks.
	Pipelines []Pipeline `json:"pipelines,omitempty"`
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

// Pipeline represents a Tekton pipeline run.
type Pipeline struct {
	Triggers     []Trigger                    `json:"triggers,omitempty"`
	Tasks        []tekton.PipelineTask        `json:"tasks,omitempty"`
	Finally      []tekton.PipelineTask        `json:"finally,omitempty"`
	Timeouts     *tekton.TimeoutFields        `json:"timeouts,omitempty"`
	PodTemplate  *pod.PodTemplate             `json:"podTemplate,omitempty"`
	TaskRunSpecs []tekton.PipelineTaskRunSpec `json:"taskRunSpecs,omitempty"`
}

// Trigger connects an incoming event with a pipeline.
type Trigger struct {
	Events         []string       `json:"events,omitempty"`
	Branches       []string       `json:"branches,omitempty"`
	ExceptBranches []string       `json:"exceptBranches,omitempty"`
	PrComment      *string        `json:"prComment,omitempty"`
	Pipeline       string         `json:"pipeline,omitempty"`
	Params         []tekton.Param `json:"params,omitempty"`
}

// Read reads an ods config from given byte slice or errors.
func Read(body []byte) (*ODS, error) {
	if len(body) == 0 {
		return nil, errors.New("config is empty")
	}
	var odsConfig *ODS
	err := yaml.UnmarshalStrict(body, &odsConfig, func(dec *json.Decoder) *json.Decoder {
		dec.DisallowUnknownFields()
		return dec
	})
	return odsConfig, err
}

// ReadFromFile reads an ods config from given filename or errors.
func ReadFromFile(filename string) (*ODS, error) {
	body, err := os.ReadFile(filename)
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
