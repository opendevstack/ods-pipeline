package tasktesting

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/api/core/v1"
)

type TestOpts struct {
	TaskKindRef             string
	TaskName                string
	Clients                 *kubernetes.Clients
	Namespace               string
	Timeout                 time.Duration
	AlwaysKeepTmpWorkspaces bool
}

type TestCase struct {
	// TaskVariant allows to target a variant of the task. E.g. the `ods-build-go task`
	// has a `ods-build-go-with-sidecar` variant. This variant can be targeted by
	// setting TaskVariant to `with-sidecar`.
	TaskVariant string
	// Map workspace name of task to local directory under test/testdata/workspaces.
	WorkspaceDirMapping map[string]string
	WantRunSuccess      bool
	PreRunFunc          func(t *testing.T, ctxt *TaskRunContext)
	PostRunFunc         func(t *testing.T, ctxt *TaskRunContext)
	Timeout             time.Duration
}

type TaskRunContext struct {
	Namespace     string
	Clients       *kubernetes.Clients
	Workspaces    map[string]string
	Params        map[string]string
	ODS           *pipelinectxt.ODSContext
	Cleanup       func()
	CollectedLogs []byte
}

func Run(t *testing.T, tc TestCase, testOpts TestOpts) {

	// Set default timeout for running the test
	if testOpts.Timeout == 0 {
		testOpts.Timeout = 120 * time.Second
	}

	taskWorkspaces := map[string]string{}
	for wn, wd := range tc.WorkspaceDirMapping {
		tempDir, err := InitWorkspace(wn, wd)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Workspace is in %s", tempDir)
		taskWorkspaces[wn] = tempDir
	}

	testCaseContext := &TaskRunContext{
		Namespace:  testOpts.Namespace,
		Clients:    testOpts.Clients,
		Workspaces: taskWorkspaces,
	}

	if tc.PreRunFunc != nil {
		tc.PreRunFunc(t, testCaseContext)
	}

	if testCaseContext.Cleanup != nil {
		defer testCaseContext.Cleanup()
	}

	tr, err := CreateTaskRunWithParams(
		testOpts.Clients.TektonClientSet,
		testOpts.TaskKindRef,
		testOpts.TaskName,
		testCaseContext.Params,
		taskWorkspaces,
		testOpts.Namespace,
	)
	if err != nil {
		t.Fatal(err)
	}

	taskRun, collectedLogsBuffer, err := WatchTaskRunUntilDone(t, testOpts, tr)
	if err != nil {
		t.Fatal(err)
	}

	if collectedLogsBuffer.Len() > 0 {
		testCaseContext.CollectedLogs = collectedLogsBuffer.Bytes()
	}

	// Show info from Task result
	CollectTaskResultInfo(taskRun, t.Logf)

	// Check if task was successful
	if taskRun.IsSuccessful() != tc.WantRunSuccess {
		t.Fatalf("Got: %+v, want: %+v.", taskRun.IsSuccessful(), tc.WantRunSuccess)
	}

	// Check local folder and evaluate output of task if needed
	if tc.PostRunFunc != nil {
		tc.PostRunFunc(t, testCaseContext)
	}

	if !testOpts.AlwaysKeepTmpWorkspaces {
		// Clean up only if test is successful
		for _, wd := range taskWorkspaces {
			err = os.RemoveAll(wd)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func InitWorkspace(workspaceName, workspaceDir string) (string, error) {
	workspaceSourceDirectory := filepath.Join(
		projectpath.Root, "test", TestdataWorkspacesPath, workspaceDir,
	)
	workspaceParentDirectory := filepath.Dir(workspaceSourceDirectory)
	return directory.CopyToTempDir(
		workspaceSourceDirectory,
		workspaceParentDirectory,
		"workspace-",
	)
}

func WatchTaskRunUntilDone(t *testing.T, testOpts TestOpts, tr *tekton.TaskRun) (*tekton.TaskRun, bytes.Buffer, error) {
	taskRunDone := make(chan *tekton.TaskRun)
	podAdded := make(chan *v1.Pod)
	errs := make(chan error)
	collectedLogsChan := make(chan []byte)
	var collectedLogsBuffer bytes.Buffer

	ctx, cancel := context.WithTimeout(context.TODO(), testOpts.Timeout)
	go waitForTaskRunDone(
		ctx,
		t,
		testOpts.Clients.TektonClientSet,
		tr.Name,
		testOpts.Namespace,
		errs,
		taskRunDone,
	)

	go waitForTaskRunPod(
		ctx,
		testOpts.Clients.KubernetesClientSet,
		tr.Name,
		testOpts.Namespace,
		podAdded,
	)

	for {
		select {
		case err := <-errs:
			if err != nil {
				cancel()
				return nil, collectedLogsBuffer, err
			}

		case pod := <-podAdded:
			if pod != nil {
				go getEventsAndLogsOfPod(
					ctx,
					testOpts.Clients.KubernetesClientSet,
					pod,
					collectedLogsChan,
					errs,
				)
			}

		case b := <-collectedLogsChan:
			collectedLogsBuffer.Write(b)

		case tr := <-taskRunDone:
			cancel()
			return tr, collectedLogsBuffer, nil
		case <-ctx.Done():
			cancel()
			return nil, collectedLogsBuffer, fmt.Errorf("timeout waiting for task run to finish. Consider increasing the timeout for your testcase at hand")
		}
	}
}
