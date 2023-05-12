package artifact

type PipelineRun struct {
	// Name is the pipeline run name.
	Name string `json:"name"`
	// AggregateTaskStatus is the aggregate Tekton task status.
	AggregateTaskStatus string `json:"aggregateTaskStatus"`
	// Repositories records the Git commit SHAs of each checked out repository.
	Repositories map[string]string `json:"repositories"`
}
