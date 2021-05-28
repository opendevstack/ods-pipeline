package test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/test/framework"
)

func TestTaskODSBuildGo(t *testing.T) {

	c, ns := framework.Setup(t,
		framework.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard",                         // if using KinD, set it to "standard"
			TaskDir:          projectpath.Root + "/deploy/tasks", // relative dir where the Tekton Task YAML file is
		},
	)

	framework.CleanupOnInterrupt(func() { framework.TearDown(t, c, ns) }, t.Logf)
	defer framework.TearDown(t, c, ns)

	tests := map[string]framework.TestCase{
		"task should build go app": {
			WorkspaceDirMapping: map[string]string{"source": "go-sample-app"},
			Params: map[string]string{
				"sonar-skip":  "true",
				"go-image":    "localhost:5000/ods/ods-build-go:latest",
				"sonar-image": "localhost:5000/ods/ods-sonar:latest",
			},
			WantSuccess: true,
			CheckFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]

				b, _, err := runCmd(wsDir+"/docker/app-linux-amd64", []string{})
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != "Hello World" {
					t.Fatalf("Got: %+v, want: %+v.", string(b), "Hello World")
				}
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			framework.Run(t, tc, framework.TestOpts{
				TaskKindRef: "ClusterTask",       // could be read from task definition
				TaskName:    "ods-build-go-v0-1", // could be read from task definition
				Clients:     c,
				Namespace:   ns,
			})

		})

	}
}

func runCmd(executable string, args []string) (outBytes, errBytes []byte, err error) {
	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outBytes = stdout.Bytes()
	errBytes = stderr.Bytes()
	return outBytes, errBytes, err
}
