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
			"clones child repositories": {
				WorkspaceDirMapping: map[string]string{"source": "multi-repo-pipeline"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					// set up bb repo for umbrella repo
					ctxt.ODS = tasktesting.SetupBitbucketRepo(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, bitbucketProjectKey)
					// set up bb repo for child repo
					tasktesting.SetupBitbucketRepo(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir+"/go-http-server", bitbucketProjectKey)
					ctxt.Params = map[string]string{
						"image":             "localhost:5000/ods/ods-start:latest",
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"console-url":       "http://example.com",
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					// umbrella
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
					// child repo
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/component", ctxt.ODS.Component)
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/git-commit-sha", ctxt.ODS.GitCommitSHA)
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/git-full-ref", ctxt.ODS.GitFullRef)
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/git-ref", ctxt.ODS.GitRef)
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/git-url", ctxt.ODS.GitURL)
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/namespace", ctxt.Namespace)
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/pr-base", "")
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/pr-key", "")
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/project", ctxt.ODS.Project)
					checkFileContent(t, wsDir, ".ods/repos/go-http-server/repository", ctxt.ODS.Repository)
				},
			},
		},
	)
}
