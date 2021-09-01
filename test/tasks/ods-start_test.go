package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSStart(t *testing.T) {
	var subrepoContext *pipelinectxt.ODSContext
	runTaskTestCases(t,
		"ods-start",
		map[string]tasktesting.TestCase{
			"clones repo": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey,
					)

					nexusClient := tasktesting.NexusClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					groupBase := fmt.Sprintf("/%s/%s/%s", ctxt.ODS.Project, ctxt.ODS.Repository, ctxt.ODS.GitCommitSHA)
					artifactsBaseDir := filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "hello-world-app-with-artifacts", pipelinectxt.ArtifactsPath)
					// Upload artifact to permanent storage.
					err := nexusClient.Upload(
						nexus.PermanentRepositoryDefault,
						groupBase+"/pipeline-runs",
						filepath.Join(artifactsBaseDir, "pipeline-runs", "foo-zh9gt0.json"),
					)
					if err != nil {
						t.Fatal(err)
					}

					ctxt.Params = map[string]string{
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"project":           ctxt.ODS.Project,
						"environment":       ctxt.ODS.Environment,
						"version":           ctxt.ODS.Version,
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					checkODSContext(t, wsDir, ctxt.ODS)

					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusInProgress)

					downloadedArtifact := filepath.Join(wsDir, pipelinectxt.PipelineRunsPath, "foo-zh9gt0.json")
					if _, err := os.Stat(downloadedArtifact); os.IsNotExist(err) {
						t.Fatal(err)
					}
				},
			},
			"clones repo and configured subrepos": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					// Setup sub-component
					tempDir, err := directory.CopyToTempDir(
						filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "hello-world-app"),
						wsDir,
						"subcomponent-",
					)
					if err != nil {
						t.Fatal(err)
					}
					subCtxt := tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, tempDir, tasktesting.BitbucketProjectKey,
					)
					subrepoContext = subCtxt
					err = os.RemoveAll(tempDir)
					if err != nil {
						t.Fatal(err)
					}
					err = createStartODSYML(wsDir, filepath.Base(tempDir))
					if err != nil {
						t.Fatal(err)
					}
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey,
					)

					nexusClient := tasktesting.NexusClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					groupBase := fmt.Sprintf("/%s/%s/%s", subCtxt.Project, subCtxt.Repository, subCtxt.GitCommitSHA)
					artifactsBaseDir := filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "hello-world-app-with-artifacts", pipelinectxt.ArtifactsPath)
					err = nexusClient.Upload(
						nexus.TemporaryRepositoryDefault,
						groupBase+"/xunit-reports",
						filepath.Join(artifactsBaseDir, "xunit-reports", "report.xml"),
					)
					if err != nil {
						t.Fatal(err)
					}
					err = nexusClient.Upload(
						nexus.TemporaryRepositoryDefault,
						groupBase+"/pipeline-runs",
						filepath.Join(artifactsBaseDir, "pipeline-runs", "foo-zh9gt0.json"),
					)
					if err != nil {
						t.Fatal(err)
					}

					ctxt.Params = map[string]string{
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"project":           ctxt.ODS.Project,
						"environment":       ctxt.ODS.Environment,
						"version":           ctxt.ODS.Version,
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					// Check .ods directory contents of main repo
					checkODSContext(t, wsDir, ctxt.ODS)

					// Check .ods directory contents of subrepo
					subrepoDir := filepath.Join(wsDir, pipelinectxt.SubreposPath, subrepoContext.Repository)
					checkODSContext(t, subrepoDir, subrepoContext)

					// Check artifacts are downloaded properly in subrepo
					sourceArtifactsBaseDir := filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "hello-world-app-with-artifacts", pipelinectxt.ArtifactsPath)
					xUnitFileSource := "xunit-reports/report.xml"
					xUnitContent := trimmedFileContentOrFatal(t, filepath.Join(sourceArtifactsBaseDir, xUnitFileSource))
					destinationArtifactsBaseDir := filepath.Join(subrepoDir, pipelinectxt.ArtifactsPath)
					checkFileContent(t, destinationArtifactsBaseDir, xUnitFileSource, xUnitContent)

					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusInProgress)

				},
			},
			"fails when subrepo has no successful pipeline run": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					// Setup sub-component
					tempDir, err := directory.CopyToTempDir(
						filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "hello-world-app"),
						wsDir,
						"subcomponent-",
					)
					if err != nil {
						t.Fatal(err)
					}
					tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, tempDir, tasktesting.BitbucketProjectKey,
					)
					err = os.RemoveAll(tempDir)
					if err != nil {
						t.Fatal(err)
					}
					err = createStartODSYML(wsDir, filepath.Base(tempDir))
					if err != nil {
						t.Fatal(err)
					}
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey,
					)
					ctxt.Params = map[string]string{
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"project":           ctxt.ODS.Project,
						"environment":       ctxt.ODS.Environment,
						"version":           ctxt.ODS.Version,
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: false,
				// TODO: check in post run func that failure is actually due to
				// missing pipeline run artifact and not due to sth. else.
			},
		},
	)
}

func createStartODSYML(wsDir, repo string) error {
	o := &config.ODS{
		Repositories: []config.Repository{
			{
				Name: repo,
			},
		},
	}
	return createODSYML(wsDir, o)
}
