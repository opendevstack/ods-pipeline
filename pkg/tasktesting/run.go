package tasktesting

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

const (
	namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
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
	// Map workspace name of task to local directory under test/testdata/workspaces.
	WorkspaceDirMapping map[string]string
	WantRunSuccess      bool
	PreRunFunc          func(t *testing.T, ctxt *TaskRunContext)
	PostRunFunc         func(t *testing.T, ctxt *TaskRunContext)
}

type TaskRunContext struct {
	Workspaces map[string]string
	Params     map[string]string
	ODS        *pipelinectxt.ODSContext
}

func Run(t *testing.T, tc TestCase, testOpts TestOpts) {

	// Set default timeout for running the test
	if testOpts.Timeout == 0 {
		testOpts.Timeout = 120 * time.Second
	}

	taskWorkspaces := map[string]string{}
	for wn, wd := range tc.WorkspaceDirMapping {
		workspaceSourceDirectory := filepath.Join(
			projectpath.Root, "test", testdataWorkspacePath, wd,
		)

		workspaceParentDirectory := filepath.Dir(workspaceSourceDirectory)

		tempDir, err := ioutil.TempDir(workspaceParentDirectory, "workspace-")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Workspace is in %s", tempDir)
		taskWorkspaces[wn] = tempDir

		directory.Copy(workspaceSourceDirectory, tempDir)

	}

	testCaseContext := &TaskRunContext{
		Workspaces: taskWorkspaces,
	}

	if tc.PreRunFunc != nil {
		tc.PreRunFunc(t, testCaseContext)
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

	WatchTaskRunEvents(testOpts.Clients.KubernetesClientSet, tr.Name, testOpts.Namespace)

	// Wait X minutes for task to complete.
	tr = WaitForCondition(context.TODO(), t, testOpts.Clients.TektonClientSet, tr.Name, testOpts.Namespace, Done, testOpts.Timeout)

	// Show logs
	CollectPodLogs(testOpts.Clients.KubernetesClientSet, tr.Status.PodName, testOpts.Namespace, t.Logf)

	// Show info from Task result
	CollectTaskResultInfo(tr, t.Logf)

	// Check if task was successful
	if tr.IsSuccessful() != tc.WantRunSuccess {
		t.Fatalf("Got: %+v, want: %+v.", tr.IsSuccessful(), tc.WantRunSuccess)
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
