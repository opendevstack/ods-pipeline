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
	Name string `json:"name"`
	URL  string `json:"url"`
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
	// Name of the environment to deploy to. This is an arbitary name.
	Name string `json:"name"`
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

type Phases struct {
	Init     Phase `json:"init"`
	Build    Phase `json:"build"`
	Deploy   Phase `json:"deploy"`
	Test     Phase `json:"test"`
	Release  Phase `json:"release"`
	Finalize Phase `json:"finalize"`
}

type Phase struct {
	RunPolicy string                `json:"runPolicy"`
	Tasks     []tekton.PipelineTask `json:"tasks"`
}
