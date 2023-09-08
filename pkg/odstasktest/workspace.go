package odstasktest

import (
	"testing"

	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/ods-pipeline/pkg/tasktesting"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
)

// GetSourceWorkspaceContext reads the ODS context from the source workspace.
func GetSourceWorkspaceContext(t *testing.T, config *ttr.TaskRunConfig) (dir string, ctxt *pipelinectxt.ODSContext) {
	dir = config.WorkspaceConfigs["source"].Dir
	ctxt, err := pipelinectxt.NewFromCache(dir)
	if err != nil {
		t.Fatal(err)
	}
	return
}

// InitGitRepo initialises a Git repository inside the given workspace.
// The workspace will also be setup with an ODS context directory in .ods
// with the given namespace.
func InitGitRepo(t *testing.T, namespace string) ttr.WorkspaceOpt {
	return func(c *ttr.WorkspaceConfig) error {
		_ = tasktesting.SetupGitRepo(t, namespace, c.Dir)
		return nil
	}
}

// WithGitSourceWorkspace configures the task run with a workspace named
// "source", mapped to the directory sourced from sourceDir. The directory is
// initialised as a Git repository with an ODS context with the given namespace.
func WithGitSourceWorkspace(t *testing.T, sourceDir, namespace string, opts ...ttr.WorkspaceOpt) ttr.TaskRunOpt {
	return WithSourceWorkspace(
		t, sourceDir,
		append([]ttr.WorkspaceOpt{InitGitRepo(t, namespace)}, opts...)...,
	)
}

// WithSourceWorkspace configures the task run with a workspace named
// "source", mapped to the directory sourced from sourceDir.
func WithSourceWorkspace(t *testing.T, sourceDir string, opts ...ttr.WorkspaceOpt) ttr.TaskRunOpt {
	return ttr.WithWorkspace("source", sourceDir, opts...)
}
