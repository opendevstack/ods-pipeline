package config

import (
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

type ODS struct {
	Repositories []Repository `json:"repositories"`
	Environments Environments `json:"environments"`

	Phases Phases `json:"phases"`
}

type Repository struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
}

type Environments struct {
	DEV  Environment `json:"dev"`
	QA   Environment `json:"qa"`
	PROD Environment `json:"prod"`
}
type Environment struct {
	Targets []Target `json:"targets"`
}

type Target struct {
	Name              string                 `json:"name"`
	URL               string                 `json:"url"`
	RegistryHost      string                 `json:"registryHost"`
	RegistryTLSVerify *bool                  `json:"registryVerify"`
	Namespace         string                 `json:"namespace"`
	SecretRef         string                 `json:"secretRef"`
	Config            map[string]interface{} `json:"config"`
}

type Phases struct {
	Init     []tekton.PipelineTask `json:"init"`
	Build    []tekton.PipelineTask `json:"build"`
	Deploy   []tekton.PipelineTask `json:"deploy"`
	Test     []tekton.PipelineTask `json:"test"`
	Release  []tekton.PipelineTask `json:"release"`
	Finalize []tekton.PipelineTask `json:"finalize"`
}
