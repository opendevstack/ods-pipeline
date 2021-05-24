package test

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/test/framework"
)

func TestTaskHelloWorld(t *testing.T) {

	taskName := "hello-world"
	workspaceName := "source" // must exist in the Task definition

	c, ns := framework.Setup(t,
		framework.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard",                           // if using KinD, set it to "standard"
			TaskDir:          "../../../../deploy/hello-world/1.0", // relative dir where the Tekton Task YAML file is
		},
	)

	framework.CleanupOnInterrupt(func() { framework.TearDown(t, c, ns) }, t.Logf)
	defer framework.TearDown(t, c, ns)

	tests := map[string]struct {
		params          map[string]string
		wantSuccess     bool
		wantFileContent string
	}{
		"task output should print Hello World": {
			params:          map[string]string{"message": "World"},
			wantSuccess:     true,
			wantFileContent: "Hello World",
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {
			tr, err := framework.CreateTaskRunWithParams(c.TektonClientSet, taskName, tc.params, workspaceName, ns)
			if err != nil {
				t.Fatal(err)
			}

			// Wait 2 minutes for task to complete.
			tr = framework.WaitForCondition(context.TODO(), t, c.TektonClientSet, tr.Name, ns, framework.Done, 120*time.Second)

			// Show logs
			framework.CollectPodLogs(c.KubernetesClientSet, tr.Status.PodName, ns, t.Logf)

			// Show info from Task result
			framework.CollectTaskResultInfo(tr, t.Logf)

			// Check if task was successful
			if tr.IsSuccessful() != tc.wantSuccess {
				t.Errorf("Got: %+v, want: %+v.", tr.IsSuccessful(), tc.wantSuccess)
			}

			// Check local folder and evaluate output of task if needed
			content, err := ioutil.ReadFile("../../../" + "msg.txt")
			if err != nil {
				t.Fatal(err)
			}

			if string(content) != tc.wantFileContent {
				t.Errorf("Got: %+v, want: %+v.", string(content), tc.wantFileContent)
			}

		})

	}
}
