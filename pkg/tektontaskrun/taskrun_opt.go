package tektontaskrun

import (
	"bytes"
	"errors"
	"log"
	"os"
	"time"

	"github.com/opendevstack/ods-pipeline/internal/directory"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

const (
	DefaultTimeout = 5 * time.Minute
)

// TaskRunOpt allows to configure the Tekton task run before it is started.
type TaskRunOpt func(c *TaskRunConfig) error

// TaskRunConfig represents key configuration of the Tekton task run.
type TaskRunConfig struct {
	Name               string
	Params             []tekton.Param
	Workspaces         map[string]string
	Namespace          string
	ServiceAccountName string
	Timeout            time.Duration
	AfterRunFunc       func(config *TaskRunConfig, taskRun *tekton.TaskRun, logs bytes.Buffer)
	CleanupFuncs       []func()
	NamespaceConfig    *NamespaceConfig
	WorkspaceConfigs   map[string]*WorkspaceConfig
	ExpectFailure      bool
}

// Cleanup calls all registered CleanupFuncs.
func (nc *TaskRunConfig) Cleanup() {
	for _, f := range nc.CleanupFuncs {
		f()
	}
}

// RunTask executes a task run after applying all given TaskRunOpt.
func RunTask(opts ...TaskRunOpt) error {
	trc := &TaskRunConfig{
		Workspaces:         map[string]string{},
		WorkspaceConfigs:   map[string]*WorkspaceConfig{},
		Timeout:            DefaultTimeout,
		ServiceAccountName: DefaultServiceAccountName,
	}
	for _, o := range opts {
		err := o(trc)
		if err != nil {
			return err
		}
	}

	cleanupOnInterrupt(trc.Cleanup)
	defer trc.Cleanup()

	taskRun, logsBuffer, err := runTask(trc)
	if err != nil {
		return err
	}

	if !taskRun.IsSuccessful() && !trc.ExpectFailure {
		return errors.New("task run was not successful")
	}

	if trc.AfterRunFunc != nil {
		trc.AfterRunFunc(trc, taskRun, logsBuffer)
	}

	return err
}

// InNamespace configures the task run to execute in given namespace.
func InNamespace(namespace string) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		c.Namespace = namespace
		return nil
	}
}

// InTempNamespace configures the task run to execute in a newly created,
// temporary namespace.
func InTempNamespace(cc *ClusterConfig, opts ...NamespaceOpt) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		nc, cleanup, err := SetupTempNamespace(cc, opts...)
		if err != nil {
			return err
		}
		c.Namespace = nc.Name
		c.NamespaceConfig = nc
		c.CleanupFuncs = append(c.CleanupFuncs, cleanup)
		return nil
	}
}

// UsingTask configures the task run to execute the Task identified by name in
// the configured namespace.
func UsingTask(name string) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		c.Name = name
		return nil
	}
}

// WithServiceAccountName configures the task run to execute under the
// specified serviceaccount name.
func WithServiceAccountName(name string) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		c.ServiceAccountName = name
		return nil
	}
}

// WithTimeout configures the task run to execute within the given duration.
func WithTimeout(timeout time.Duration) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		c.Timeout = timeout
		return nil
	}
}

// WithWorkspace sets up a workspace with given name and contents of sourceDir.
// sourceDir is copied to a temporary directory so that the original contents
// remain unchanged.
func WithWorkspace(name, sourceDir string, opts ...WorkspaceOpt) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		workspaceDir, cleanup, err := SetupWorkspaceDir(sourceDir)
		if err != nil {
			return err
		}
		log.Printf("Workspace %q is in %s ...\n", name, workspaceDir)
		wc := &WorkspaceConfig{
			Name:    name,
			Dir:     workspaceDir,
			Cleanup: cleanup,
		}
		for _, o := range opts {
			err := o(wc)
			if err != nil {
				return err
			}
		}
		c.WorkspaceConfigs[wc.Name] = wc
		c.CleanupFuncs = append(c.CleanupFuncs, wc.Cleanup)
		c.Workspaces[wc.Name] = wc.Dir
		return nil
	}
}

// WithParams configures the task run to use the specified Tekton parameters.
func WithParams(params ...tekton.Param) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		c.Params = append(c.Params, params...)
		return nil
	}
}

// WithStringParams configures the task run to use the specified string
// parameters. WithStringParams is a more convenient way to configure
// simple parameters compares to WithParams.
func WithStringParams(params map[string]string) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		c.Params = append(c.Params, TektonParamsFromStringParams(params)...)
		return nil
	}
}

// ExpectFailure sets up an expectation that the task will fail. If the task
// does not fail, RunTask will error. Conversely, if ExpectFailure is not set,
// RunTask will error when the task run fails.
func ExpectFailure() TaskRunOpt {
	return func(c *TaskRunConfig) error {
		c.ExpectFailure = true
		return nil
	}
}

// AfterRun registers a function which is run after the task run completes.
// The function will receive the task run configuration, as well as an instance
// of the TaskRun.
func AfterRun(f func(c *TaskRunConfig, r *tekton.TaskRun, l bytes.Buffer)) TaskRunOpt {
	return func(c *TaskRunConfig) error {
		c.AfterRunFunc = f
		return nil
	}
}

// SetupWorkspaceDir copies sourceDir to the KinD mount host path, which is
// set to /tmp/ods-pipeline/kind-mount. The created folder can then be used
// as a Tekton task run workspace. SetupWorkspaceDir returns the
// created directory as well as a function to clean it up.
func SetupWorkspaceDir(sourceDir string) (dir string, cleanup func(), err error) {
	dir, err = directory.CopyToTempDir(sourceDir, KinDMountHostPath, "workspace-")
	cleanup = func() {
		if err := os.RemoveAll(dir); err != nil {
			log.Printf("failed to clean up temporary workspace dir %s: %s", dir, err)
		}
	}
	return
}
