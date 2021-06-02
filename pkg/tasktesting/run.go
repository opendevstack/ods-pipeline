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
	Params              map[string]string
	WantSuccess         bool
	PrepareFunc         func(t *testing.T, workspaces map[string]string)
	CheckFunc           func(t *testing.T, workspaces map[string]string)
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

	if tc.PrepareFunc != nil {
		tc.PrepareFunc(t, taskWorkspaces)
	}

	tr, err := CreateTaskRunWithParams(
		testOpts.Clients.TektonClientSet,
		testOpts.TaskKindRef,
		testOpts.TaskName,
		tc.Params,
		taskWorkspaces,
		testOpts.Namespace,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Wait X minutes for task to complete.
	tr = WaitForCondition(context.TODO(), t, testOpts.Clients.TektonClientSet, tr.Name, testOpts.Namespace, Done, testOpts.Timeout)

	// Show logs
	CollectPodLogs(testOpts.Clients.KubernetesClientSet, tr.Status.PodName, testOpts.Namespace, t.Logf)

	// Show info from Task result
	CollectTaskResultInfo(tr, t.Logf)

	// Check if task was successful
	if tr.IsSuccessful() != tc.WantSuccess {
		t.Fatalf("Got: %+v, want: %+v.", tr.IsSuccessful(), tc.WantSuccess)
	}

	// Check local folder and evaluate output of task if needed
	tc.CheckFunc(t, taskWorkspaces)

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
