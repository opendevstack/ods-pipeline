package tasks

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSStart(t *testing.T) {

	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard", // if using KinD, set it to "standard"
			TaskDir:          projectpath.Root + "/deploy/tasks",
			EnvironmentDir:   projectpath.Root + "/test/testdata/deploy/cd-kind",
		},
	)

	tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	defer tasktesting.TearDown(t, c, ns)

	var originURL string
	tests := map[string]tasktesting.TestCase{
		"clones the app": {
			WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
			Params: map[string]string{
				"image":      "localhost:5000/ods/start:latest",
				"url":        "",     // set later in prepare func
				"project":    "proj", // TODO: make this param optional?
				"component":  "comp", // TODO: make this param optional?
				"repository": "repo", // TODO: make this param optional?
			},
			PrepareFunc: func(t *testing.T, workspaces, params map[string]string) {
				wsDir := workspaces["source"]
				os.Chdir(wsDir)
				tasktesting.InitAndCommitOrFatal(t, wsDir) // will be cleaned by task
				originURL = tasktesting.PushToBitbucketOrFatal(t, c.KubernetesClientSet, ns, wsDir, bitbucketProjectKey)
				params["url"] = originURL
			},
			WantSuccess: true,
			CheckFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]

				wantFiles := []string{
					"docker/Dockerfile",
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

			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             "ClusterTask",      // could be read from task definition
				TaskName:                "ods-start-v0-1-0", // could be read from task definition
				Clients:                 c,
				Namespace:               ns,
				Timeout:                 5 * time.Minute, // depending on  the task we may need to increase or decrease it
				AlwaysKeepTmpWorkspaces: *alwaysKeepTmpWorkspacesFlag,
			})

		})

	}
}
