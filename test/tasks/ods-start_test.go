package tasks

import (
	"os"
	"path/filepath"
	"strings"
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
	var lfsFilename string
	var lfsFileHash [32]byte
	runTaskTestCases(t,
		"ods-start",
		[]tasktesting.Service{
			tasktesting.Bitbucket,
			tasktesting.Nexus,
		},
		map[string]tasktesting.TestCase{
			"clones repo": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
					)
					ctxt.Params = map[string]string{
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"project":           ctxt.ODS.Project,
						"version":           ctxt.ODS.Version,
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					checkODSContext(t, wsDir, ctxt.ODS)

					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, *privateCertFlag)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusInProgress)
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
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, tempDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
					)
					subrepoContext = subCtxt
					err = os.RemoveAll(tempDir)
					if err != nil {
						t.Fatal(err)
					}
					err = createStartODSYMLWithSubrepo(wsDir, filepath.Base(tempDir))
					if err != nil {
						t.Fatal(err)
					}
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
					)

					nexusClient := tasktesting.NexusClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, *privateCertFlag)
					artifactsBaseDir := filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "hello-world-app-with-artifacts", pipelinectxt.ArtifactsPath)
					_, err = nexusClient.Upload(
						nexus.TestTemporaryRepository,
						pipelinectxt.ArtifactGroup(subCtxt, pipelinectxt.XUnitReportsDir),
						filepath.Join(artifactsBaseDir, pipelinectxt.XUnitReportsDir, "report.xml"),
					)
					if err != nil {
						t.Fatal(err)
					}
					_, err = nexusClient.Upload(
						nexus.TestTemporaryRepository,
						pipelinectxt.ArtifactGroup(subCtxt, pipelinectxt.PipelineRunsDir),
						filepath.Join(artifactsBaseDir, pipelinectxt.PipelineRunsDir, "foo-zh9gt0.json"),
					)
					if err != nil {
						t.Fatal(err)
					}

					ctxt.Params = map[string]string{
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"project":           ctxt.ODS.Project,
						"version":           ctxt.ODS.Version,
						"pipeline-run-name": "foo",
						"artifact-source":   nexus.TestTemporaryRepository,
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					// Check .ods directory contents of main repo
					checkODSContext(t, wsDir, ctxt.ODS)
					checkFilesExist(t, wsDir, filepath.Join(pipelinectxt.ArtifactsPath, pipelinectxt.ArtifactsManifestFilename))

					// Check .ods directory contents of subrepo
					subrepoDir := filepath.Join(wsDir, pipelinectxt.SubreposPath, subrepoContext.Repository)
					checkODSContext(t, subrepoDir, subrepoContext)

					// Check artifacts are downloaded properly in subrepo
					sourceArtifactsBaseDir := filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "hello-world-app-with-artifacts", pipelinectxt.ArtifactsPath)
					xUnitFileSource := "xunit-reports/report.xml"
					xUnitContent := trimmedFileContentOrFatal(t, filepath.Join(sourceArtifactsBaseDir, xUnitFileSource))
					destinationArtifactsBaseDir := filepath.Join(subrepoDir, pipelinectxt.ArtifactsPath)
					checkFileContent(t, destinationArtifactsBaseDir, xUnitFileSource, xUnitContent)
					checkFilesExist(t, destinationArtifactsBaseDir, pipelinectxt.ArtifactsManifestFilename)

					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, *privateCertFlag)
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
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, tempDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
					)
					err = os.RemoveAll(tempDir)
					if err != nil {
						t.Fatal(err)
					}
					err = createStartODSYMLWithSubrepo(wsDir, filepath.Base(tempDir))
					if err != nil {
						t.Fatal(err)
					}
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
					)
					ctxt.Params = map[string]string{
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"project":           ctxt.ODS.Project,
						"version":           ctxt.ODS.Version,
						"pipeline-run-name": "foo",
						"artifact-source":   "empty-repo",
					}
				},
				WantRunSuccess: false,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					want := "Pipeline runs with subrepos require a successful pipeline run artifact " +
						"for all checked out subrepo commits, however no such artifact was found"

					if !strings.Contains(string(ctxt.CollectedLogs), want) {
						t.Fatalf("Want:\n%s\n\nGot:\n%s", want, string(ctxt.CollectedLogs))
					}
				},
			},
			"handles git LFS extension": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
					)
					tasktesting.EnableLfsOnBitbucketRepoOrFatal(t, filepath.Base(wsDir), tasktesting.BitbucketProjectKey)
					lfsFilename = "lfspicture.jpg"
					lfsFileHash = tasktesting.UpdateBitbucketRepoWithLfsOrFatal(t, ctxt.ODS, wsDir, tasktesting.BitbucketProjectKey, lfsFilename)

					ctxt.Params = map[string]string{
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"project":           ctxt.ODS.Project,
						"version":           ctxt.ODS.Version,
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					checkODSContext(t, wsDir, ctxt.ODS)

					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, *privateCertFlag)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusInProgress)

					checkFileHash(t, wsDir, lfsFilename, lfsFileHash)
				},
			},
		},
	)
}

func createStartODSYMLWithSubrepo(wsDir, repo string) error {
	o := &config.ODS{
		Repositories: []config.Repository{
			{
				Name: repo,
			},
		},
	}
	return createODSYML(wsDir, o)
}
