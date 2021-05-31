package tasks

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/test/framework"
)

func TestTaskODSBuildImage(t *testing.T) {

	c, ns := framework.Setup(t,
		framework.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard", // if using KinD, set it to "standard"
			TaskDir:          projectpath.Root + "/deploy/tasks",
			EnvironmentDir:   projectpath.Root + "/test/testdata/deploy/cd-kind",
		},
	)

	framework.CleanupOnInterrupt(func() { framework.TearDown(t, c, ns) }, t.Logf)
	defer framework.TearDown(t, c, ns)

	tests := map[string]framework.TestCase{
		"task should build image": {
			WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
			Params: map[string]string{
				"registry":      "kind-registry.kind:5000",
				"builder-image": "localhost:5000/ods/ods-buildah:latest",
				"tls-verify":    "false",
			},
			WantSuccess: true,
			CheckFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]

				wantFiles := []string{
					".ods/artifacts/image-digests/hello-world.json",
				}
				for _, wf := range wantFiles {
					if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
						t.Fatalf("Want %s, but got nothing", wf)
					}
				}
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			framework.Run(t, tc, framework.TestOpts{
				TaskKindRef: "ClusterTask",          // could be read from task definition
				TaskName:    "ods-build-image-v0-1", // could be read from task definition
				Clients:     c,
				Namespace:   ns,
				Timeout:     5 * time.Minute, // depending on  the task we may need to increase or decrease it
			})

		})

	}
}
