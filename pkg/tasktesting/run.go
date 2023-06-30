package tasktesting

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/ods-pipeline/internal/kubernetes"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
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

type TaskRunCase struct {
	WantRunSuccess bool
	WantSetupFail  bool
	PreRunFunc     func(t *testing.T, ctxt *TaskRunContext)
	PostRunFunc    func(t *testing.T, ctxt *TaskRunContext)
}

type TestCase struct {
	// Map workspace name of task to local directory under test/testdata/workspaces.
	WorkspaceDirMapping map[string]string
	TaskParamsMapping   map[string]string
	// The fields until AdditionalRuns are the same as for AdditionalRuns TaskRunCase so that single task tests do not require another level. The goal is to only require the extra level if needed. Could this be avoided?
	WantRunSuccess bool
	WantSetupFail  bool
	PreRunFunc     func(t *testing.T, ctxt *TaskRunContext)
	PostRunFunc    func(t *testing.T, ctxt *TaskRunContext)
	CleanupFunc    func(t *testing.T, ctxt *TaskRunContext)
	AdditionalRuns []TaskRunCase
	Timeout        time.Duration
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

func (tc *TaskRunContext) toTektonParams() []tekton.Param {
	var tektonParams []tekton.Param

	// When tekton supports array usage their usage could be supported here.
	// (see [Cannot refer array params in script #4912](https://github.com/tektoncd/pipeline/issues/4912))
	for key, value := range tc.Params {
		tektonParams = append(tektonParams, tekton.Param{
			Name:  key,
			Value: *tekton.NewArrayOrString(value),
		})
	}
	return tektonParams
}

func runTask(t *testing.T, testOpts TestOpts, taskWorkspaces map[string]string, testCaseContext *TaskRunContext, tc TaskRunCase) {
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
		testCaseContext.toTektonParams(),
		taskWorkspaces,
		testOpts.Namespace,
	)
	if err != nil {
		t.Fatal(err)
	}

	taskRun, collectedLogsBuffer, err := WatchTaskRunUntilDone(t, testOpts, tr)

	// Check if task setup was successful
	if err != nil {
		if tc.WantSetupFail {
			return
		} else {
			t.Fatalf("Task setup failed: %s", err)
		}
	}

	if tc.WantSetupFail {
		t.Fatal("Task setup was successful, but was expected to fail.")
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
}

func Run(t *testing.T, tc TestCase, testOpts TestOpts) {

	// Set default timeout for running the test
	if testOpts.Timeout == 0 {
		testOpts.Timeout = 120 * time.Second
	}

	taskWorkspaces := map[string]string{}
	var workspaceCleanupFuncs []func()
	for wn, wd := range tc.WorkspaceDirMapping {
		workspaceSourceDirectory := filepath.Join(
			projectpath.Root, "test", TestdataWorkspacesPath, wd,
		)
		tempDir, cf, err := tektontaskrun.SetupWorkspaceDir(workspaceSourceDirectory)
		if err != nil {
			t.Fatal(err)
		}
		workspaceCleanupFuncs = append(workspaceCleanupFuncs, cf)
		t.Logf("Workspace is in %s", tempDir)
		taskWorkspaces[wn] = tempDir
	}

	testCaseContext := &TaskRunContext{
		Namespace:  testOpts.Namespace,
		Clients:    testOpts.Clients,
		Workspaces: taskWorkspaces,
		Params:     tc.TaskParamsMapping,
	}

	if tc.CleanupFunc != nil {
		defer tc.CleanupFunc(t, testCaseContext)
	}

	tasks := []TaskRunCase{}
	tasks = append(tasks, TaskRunCase{
		WantRunSuccess: tc.WantRunSuccess,
		WantSetupFail:  tc.WantSetupFail,
		PreRunFunc:     tc.PreRunFunc,
		PostRunFunc:    tc.PostRunFunc,
	})
	for _, ttc := range tc.AdditionalRuns {
		if ttc.PostRunFunc == nil {
			ttc.PostRunFunc = tc.PostRunFunc
		}
		if ttc.PreRunFunc == nil {
			ttc.PreRunFunc = tc.PreRunFunc
		}
		tasks = append(tasks, ttc)
	}
	for _, ttc := range tasks {
		runTask(t, testOpts, taskWorkspaces, testCaseContext, ttc)
	}

	if !testOpts.AlwaysKeepTmpWorkspaces {
		// Clean up only if test is successful
		for _, cf := range workspaceCleanupFuncs {
			cf()
		}
	}
}

func WatchTaskRunUntilDone(t *testing.T, testOpts TestOpts, tr *tekton.TaskRun) (*tekton.TaskRun, bytes.Buffer, error) {
	taskRunDone := make(chan *tekton.TaskRun)
	podAdded := make(chan *v1.Pod)
	errs := make(chan error)
	collectedLogsChan := make(chan []byte)
	var collectedLogsBuffer bytes.Buffer

	ctx, cancel := context.WithTimeout(context.TODO(), testOpts.Timeout)
	defer cancel()
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
			return tr, collectedLogsBuffer, nil
		case <-ctx.Done():
			return nil, collectedLogsBuffer, fmt.Errorf("timeout waiting for task run to finish. Consider increasing the timeout for your testcase at hand")
		}
	}
}
