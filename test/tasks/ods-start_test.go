package tasks

import (
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

	tests := map[string]tasktesting.TestCase{
		"clones the app": {
			WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
			PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				ctxt.ODS = tasktesting.SetupBitbucketRepo(t, c.KubernetesClientSet, ns, wsDir, bitbucketProjectKey)
				ctxt.Params = map[string]string{
					"image":             "localhost:5000/ods/start:latest",
					"url":               ctxt.ODS.GitURL,
					"git-full-ref":      "refs/heads/master",
					"project":           ctxt.ODS.Project,
					"component":         ctxt.ODS.Component,
					"repository":        ctxt.ODS.Repository,
					"console-url":       "http://example.com",
					"pipeline-run-name": "foo",
				}
			},
			WantRunSuccess: true,
			PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]

				checkFileContent(t, wsDir, ".ods/component", ctxt.ODS.Component)
				checkFileContent(t, wsDir, ".ods/git-commit-sha", ctxt.ODS.GitCommitSHA)
				checkFileContent(t, wsDir, ".ods/git-full-ref", ctxt.ODS.GitFullRef)
				checkFileContent(t, wsDir, ".ods/git-ref", ctxt.ODS.GitRef)
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
