package tasks

import (
	"fmt"
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

					nexusClient := tasktesting.NexusClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, *privateCertFlag)
					artifactsBaseDir := filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "hello-world-app-with-artifacts", pipelinectxt.ArtifactsPath)
					// Upload artifact to permanent storage.
					_, err := nexusClient.Upload(
						nexus.PermanentRepositoryDefault,
						pipelinectxt.ArtifactGroup(ctxt.ODS, pipelinectxt.PipelineRunsDir),
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

					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, *privateCertFlag)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusInProgress)

					checkFilesExist(t, wsDir,
						filepath.Join(pipelinectxt.PipelineRunsPath, "foo-zh9gt0.json"),
						filepath.Join(pipelinectxt.ArtifactsPath, pipelinectxt.ArtifactsManifestFilename),
					)
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
						nexus.TemporaryRepositoryDefault,
						pipelinectxt.ArtifactGroup(subCtxt, pipelinectxt.XUnitReportsDir),
						filepath.Join(artifactsBaseDir, pipelinectxt.XUnitReportsDir, "report.xml"),
					)
					if err != nil {
						t.Fatal(err)
					}
					_, err = nexusClient.Upload(
						nexus.TemporaryRepositoryDefault,
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
						"environment":       ctxt.ODS.Environment,
						"version":           ctxt.ODS.Version,
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: false,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					want := "Pipeline runs with subrepos require a successful pipeline run " +
						"for all checked out subrepo commits, however no such run was found"

					if !strings.Contains(string(ctxt.CollectedLogs), want) {
						t.Fatalf("Want:\n%s\n\nGot:\n%s", want, string(ctxt.CollectedLogs))
					}
				},
			},
			"handles QA stage": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					version := "1.0.0"
					err := createStartODSYML(wsDir, config.QAStage)
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
						"environment":       ctxt.ODS.Environment,
						"version":           version,
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					gotTag := checkForTag(t, ctxt, "v1.0.0-rc.1")
					if gotTag.LatestCommit != ctxt.ODS.GitCommitSHA {
						t.Fatal("not same SHA")
					}
				},
			},
			"handles PROD stage": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					version := "1.0.0"
					err := createStartODSYML(wsDir, config.ProdStage)
					if err != nil {
						t.Fatal(err)
					}
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
					)

					// pretend there is already an RC tag for the current commit
					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, *privateCertFlag)
					_, err = bitbucketClient.TagCreate(
						ctxt.ODS.Project,
						ctxt.ODS.Repository,
						bitbucket.TagCreatePayload{
							Name:       fmt.Sprintf("v%s-rc.1", version),
							StartPoint: ctxt.ODS.GitCommitSHA,
						},
					)
					if err != nil {
						t.Fatal(err)
					}

					ctxt.Params = map[string]string{
						"url":               ctxt.ODS.GitURL,
						"git-full-ref":      "refs/heads/master",
						"project":           ctxt.ODS.Project,
						"environment":       ctxt.ODS.Environment,
						"version":           version,
						"pipeline-run-name": "foo",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					gotTag := checkForTag(t, ctxt, "v1.0.0")
					if gotTag.LatestCommit != ctxt.ODS.GitCommitSHA {
						t.Fatalf(
							"Checked out commit is %s, but created tag points to %s",
							ctxt.ODS.GitCommitSHA,
							gotTag.LatestCommit,
						)
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
						"environment":       ctxt.ODS.Environment,
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

func checkForTag(t *testing.T, ctxt *tasktesting.TaskRunContext, wantTag string) *bitbucket.Tag {
	bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, *privateCertFlag)
	var gotTag *bitbucket.Tag
	tagPage, err := bitbucketClient.TagList(
		ctxt.ODS.Project,
		ctxt.ODS.Repository,
		bitbucket.TagListParams{
			FilterText: wantTag,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range tagPage.Values {
		if tag.DisplayID == wantTag {
			gotTag = &tag
			break
		}
	}
	if gotTag == nil {
		t.Fatalf("Could not find tag %s", wantTag)
	}
	return gotTag
}

func createStartODSYMLWithSubrepo(wsDir, repo string) error {
	o := &config.ODS{
		Environments: []config.Environment{
			{
				Name:  "dev",
				Stage: config.DevStage,
			},
		},
		Repositories: []config.Repository{
			{
				Name: repo,
			},
		},
	}
	return createODSYML(wsDir, o)
}

func createStartODSYML(wsDir string, stage config.Stage) error {
	o := &config.ODS{
		Environments: []config.Environment{
			{
				Name:  "dev",
				Stage: stage,
			},
		},
	}
	return createODSYML(wsDir, o)
}
