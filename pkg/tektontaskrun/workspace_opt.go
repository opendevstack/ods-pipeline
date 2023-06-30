package tektontaskrun

// WorkspaceOpt allows to further configure a Tekton workspace after its creation.
type WorkspaceOpt func(c *WorkspaceConfig) error

// WorkspaceConfig describes a Tekton workspace.
type WorkspaceConfig struct {
	// Name of the Tekton workspace.
	Name string
	// Directory on the host of the workspace.
	Dir string
	// Cleanup function.
	Cleanup func()
}
