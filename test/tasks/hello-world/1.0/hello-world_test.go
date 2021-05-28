package test

import (
	"io/ioutil"
	"testing"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/test/framework"
)

func TestTaskHelloWorld(t *testing.T) {

	c, ns := framework.Setup(t,
		framework.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard",                                   // if using KinD, set it to "standard"
			TaskDir:          projectpath.Root + "/deploy/hello-world/1.0", // relative dir where the Tekton Task YAML file is
		},
	)

	framework.CleanupOnInterrupt(func() { framework.TearDown(t, c, ns) }, t.Logf)
	defer framework.TearDown(t, c, ns)

	tests := map[string]framework.TestCase{
		"task output should print Hello World": {
			WorkspaceDirMapping: map[string]string{"source": "hello-world"},
			Params:              map[string]string{"message": "World"},
			WantSuccess:         true,
			CheckFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]
				content, err := ioutil.ReadFile(wsDir + "/" + "msg.txt")
				if err != nil {
					t.Fatal(err)
				}

				if string(content) != "Hello World" {
					t.Errorf("Got: %+v, want: %+v.", string(content), "Hello World")
				}
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			framework.Run(t, tc, framework.TestOpts{
				TaskKindRef: "Task",        // could be read from task definition
				TaskName:    "hello-world", // could be read from task definition
				Clients:     c,
				Namespace:   ns,
			})

		})

	}
}
