package tasks

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
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
			PreTaskRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				os.Chdir(wsDir)
				tasktesting.InitAndCommitOrFatal(t, wsDir) // will be cleaned by task
				originURL = tasktesting.PushToBitbucketOrFatal(t, c.KubernetesClientSet, ns, wsDir, bitbucketProjectKey)

				ctxt.ODS = &pipelinectxt.ODSContext{
					Namespace: ns,
					Project:   bitbucketProjectKey,
					GitURL:    originURL,
				}
				err := ctxt.ODS.Assemble(wsDir)
				if err != nil {
					t.Fatalf("could not assemble ODS context information: %s", err)
				}

				ctxt.Params = map[string]string{
					"image":      "localhost:5000/ods/start:latest",
					"url":        originURL,
					"git-ref":    "master",
					"project":    ctxt.ODS.Project,
					"component":  ctxt.ODS.Component,
					"repository": ctxt.ODS.Repository,
				}
			},
			WantSuccess: true,
			PostTaskRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]

				checkFileContent(t, wsDir, ".ods/component", ctxt.ODS.Component)
				checkFileContent(t, wsDir, ".ods/git-commit-sha", ctxt.ODS.GitCommitSHA)
				// checkFileContent(t, wsDir, ".ods/git-full-ref", ctxt.ODS.GitFullRef) // TODO: implement in task
				// checkFileContent(t, wsDir, ".ods/git-ref", ctxt.ODS.GitRef) // TODO: implement in task
				checkFileContent(t, wsDir, ".ods/git-url", ctxt.ODS.GitURL)
				checkFileContent(t, wsDir, ".ods/namespace", ns)
				checkFileContent(t, wsDir, ".ods/pr-base", "")
				checkFileContent(t, wsDir, ".ods/pr-key", "")
				checkFileContent(t, wsDir, ".ods/project", ctxt.ODS.Project)
				checkFileContent(t, wsDir, ".ods/repository", ctxt.ODS.Repository)

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

func checkFileContent(t *testing.T, wsDir, filename, want string) {
	got, err := getTrimmedFileContent(filepath.Join(wsDir, filename))
	if err != nil {
		t.Fatalf("could not read %s: %s", filename, err)
	}
	if got != want {
		t.Fatalf("got '%s', want '%s' in file %s", got, want, filename)
	}
}
