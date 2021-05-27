package framework

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
	TaskKindRef   string
	TaskName      string
	WorkspaceName string
	Clients       *kubernetes.Clients
	Namespace     string
}

type TestCase struct {
	WorkspaceSourceDirectory string
	Params                   map[string]string
	WantSuccess              bool
	CheckFunc                func(t *testing.T, workspaceDir string)
}

func Run(t *testing.T, tc TestCase, testOpts TestOpts) {

	workspaceSourceDirectory := filepath.Join(
		projectpath.Root, "test", testdataWorkspacePath, tc.WorkspaceSourceDirectory,
	)

	workspaceParentDirectory := filepath.Dir(workspaceSourceDirectory)

	tempDir, err := ioutil.TempDir(workspaceParentDirectory, "workspace-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Workspace is in %s", tempDir)
	tempDirName := filepath.Base(tempDir)

	directory.Copy(workspaceSourceDirectory, tempDir)

	tr, err := CreateTaskRunWithParams(
		testOpts.Clients.TektonClientSet,
		testOpts.TaskKindRef,
		testOpts.TaskName,
		tc.Params,
		testOpts.WorkspaceName,
		tempDirName,
		testOpts.Namespace,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Wait 2 minutes for task to complete.
	tr = WaitForCondition(context.TODO(), t, testOpts.Clients.TektonClientSet, tr.Name, testOpts.Namespace, Done, 120*time.Second)

	// Show logs
	CollectPodLogs(testOpts.Clients.KubernetesClientSet, tr.Status.PodName, testOpts.Namespace, t.Logf)

	// Show info from Task result
	CollectTaskResultInfo(tr, t.Logf)

	// Check if task was successful
	if tr.IsSuccessful() != tc.WantSuccess {
		t.Fatalf("Got: %+v, want: %+v.", tr.IsSuccessful(), tc.WantSuccess)
	}

	// Check local folder and evaluate output of task if needed
	tc.CheckFunc(t, tempDir)

	// Clean up only if test is successful
	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}
}
