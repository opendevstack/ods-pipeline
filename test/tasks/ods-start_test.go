package tasks

import (
	"testing"

	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSStart(t *testing.T) {
	runTaskTestCases(t,
		"ods-start-v0-1-0",
		map[string]tasktesting.TestCase{
			"clones the app": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupBitbucketRepo(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, bitbucketProjectKey)
					ctxt.Params = map[string]string{
						"image":             "localhost:5000/ods/ods-start:latest",
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
					checkFileContent(t, wsDir, ".ods/namespace", ctxt.Namespace)
					checkFileContent(t, wsDir, ".ods/pr-base", "")
					checkFileContent(t, wsDir, ".ods/pr-key", "")
					checkFileContent(t, wsDir, ".ods/project", ctxt.ODS.Project)
					checkFileContent(t, wsDir, ".ods/repository", ctxt.ODS.Repository)

				},
			},
		},
	)
}
