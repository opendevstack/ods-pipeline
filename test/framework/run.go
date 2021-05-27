package framework

import (
	"context"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/kubernetes"
)

type TestOpts struct {
	TaskKindRef   string
	TaskName      string
	WorkspaceName string
	Clients       *kubernetes.Clients
	Namespace     string
}

type TestCase struct {
	Params      map[string]string
	WantSuccess bool
	CheckFunc   func(t *testing.T)
}

func Run(t *testing.T, tc TestCase, testOpts TestOpts) {
	tr, err := CreateTaskRunWithParams(testOpts.Clients.TektonClientSet, testOpts.TaskKindRef, testOpts.TaskName, tc.Params, testOpts.WorkspaceName, testOpts.Namespace)
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
		t.Errorf("Got: %+v, want: %+v.", tr.IsSuccessful(), tc.WantSuccess)
	}

	// Check local folder and evaluate output of task if needed
	tc.CheckFunc(t)
}
