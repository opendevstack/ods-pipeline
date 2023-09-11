package e2e

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/ods-pipeline/internal/directory"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	"github.com/opendevstack/ods-pipeline/pkg/config"
	"github.com/opendevstack/ods-pipeline/pkg/nexus"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/ods-pipeline/pkg/tasktesting"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	ott "github.com/opendevstack/ods-pipeline/pkg/odstasktest"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
)

func runStartTask(opts ...ttr.TaskRunOpt) error {
	return ttr.RunTask(append([]ttr.TaskRunOpt{
		ttr.InNamespace(namespaceConfig.Name),
		ttr.UsingTask("ods-pipeline-start"),
	}, opts...)...)
}

func TestStartTaskClonesRepoAtBranch(t *testing.T) {
	k8sClient := newK8sClient(t)
	if err := runStartTask(
		withBitbucketSourceWorkspace(t, "../testdata/workspaces/hello-world-app", k8sClient, namespaceConfig.Name),
		func(c *ttr.TaskRunConfig) error {
			c.Params = append(c.Params, ttr.TektonParamsFromStringParams(map[string]string{
				"url":               bitbucketURLForWorkspace(c.WorkspaceConfigs["source"]),
				"git-full-ref":      "refs/heads/master",
				"project":           tasktesting.BitbucketProjectKey,
				"pipeline-run-name": "foo",
			})...)
			return nil
		},
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			wsDir, odsContext := ott.GetSourceWorkspaceContext(t, config)

			checkODSContext(t, wsDir, odsContext)
			checkFilesExist(t, wsDir, filepath.Join(pipelinectxt.ArtifactsPath, pipelinectxt.ArtifactsManifestFilename))

			bitbucketClient := tasktesting.BitbucketClientOrFatal(t, k8sClient, namespaceConfig.Name, *privateCertFlag)
			checkBuildStatus(t, bitbucketClient, odsContext.GitCommitSHA, bitbucket.BuildStatusInProgress)
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func TestStartTaskClonesRepoAtTag(t *testing.T) {
	k8sClient := newK8sClient(t)
	if err := runStartTask(
		withBitbucketSourceWorkspace(t, "../testdata/workspaces/hello-world-app", k8sClient, namespaceConfig.Name),
		func(c *ttr.TaskRunConfig) error {
			wsDir, odsContext := ott.GetSourceWorkspaceContext(t, c)
			tasktesting.UpdateBitbucketRepoWithTagOrFatal(t, odsContext, wsDir, "v1.0.0")
			c.Params = append(c.Params, ttr.TektonParamsFromStringParams(map[string]string{
				"url":               bitbucketURLForWorkspace(c.WorkspaceConfigs["source"]),
				"git-full-ref":      "refs/tags/v1.0.0",
				"project":           tasktesting.BitbucketProjectKey,
				"pipeline-run-name": "foo",
			})...)
			return nil
		},
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			wsDir, odsContext := ott.GetSourceWorkspaceContext(t, config)
			checkODSContext(t, wsDir, odsContext)
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func TestStartTaskClonesRepoAndSubrepos(t *testing.T) {
	var subrepoContext *pipelinectxt.ODSContext
	k8sClient := newK8sClient(t)
	if err := runStartTask(
		ott.WithSourceWorkspace(
			t,
			"../testdata/workspaces/hello-world-app",
			func(c *ttr.WorkspaceConfig) error {
				// Setup sub-component
				subrepoContext = setupBitbucketRepoWithSubdirOrFatal(t, c, k8sClient)
				// Nexus artifacts
				nexusClient := tasktesting.NexusClientOrFatal(t, k8sClient, namespaceConfig.Name, *privateCertFlag)
				artifactsBaseDir := filepath.Join(projectpath.Root, "test", testdataWorkspacesPath, "hello-world-app-with-artifacts", pipelinectxt.ArtifactsPath)
				_, err := nexusClient.Upload(
					nexus.TestTemporaryRepository,
					pipelinectxt.ArtifactGroup(subrepoContext, pipelinectxt.XUnitReportsDir),
					filepath.Join(artifactsBaseDir, pipelinectxt.XUnitReportsDir, "report.xml"),
				)
				if err != nil {
					return err
				}
				_, err = nexusClient.Upload(
					nexus.TestTemporaryRepository,
					pipelinectxt.ArtifactGroup(subrepoContext, pipelinectxt.PipelineRunsDir),
					filepath.Join(artifactsBaseDir, pipelinectxt.PipelineRunsDir, "foo-zh9gt0.json"),
				)
				if err != nil {
					return err
				}
				return nil
			},
		),
		func(c *ttr.TaskRunConfig) error {
			c.Params = append(c.Params, ttr.TektonParamsFromStringParams(map[string]string{
				"url":               bitbucketURLForWorkspace(c.WorkspaceConfigs["source"]),
				"git-full-ref":      "refs/heads/master",
				"project":           tasktesting.BitbucketProjectKey,
				"pipeline-run-name": "foo",
				"artifact-source":   nexus.TestTemporaryRepository,
			})...)
			return nil
		},
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			wsDir, odsContext := ott.GetSourceWorkspaceContext(t, config)

			// Check .ods directory contents of main repo
			checkODSContext(t, wsDir, odsContext)
			checkFilesExist(t, wsDir, filepath.Join(pipelinectxt.ArtifactsPath, pipelinectxt.ArtifactsManifestFilename))

			// Check .ods directory contents of subrepo
			subrepoDir := filepath.Join(wsDir, pipelinectxt.SubreposPath, subrepoContext.Repository)
			checkODSContext(t, subrepoDir, subrepoContext)

			// Check artifacts are downloaded properly in subrepo
			sourceArtifactsBaseDir := filepath.Join(projectpath.Root, "test", testdataWorkspacesPath, "hello-world-app-with-artifacts", pipelinectxt.ArtifactsPath)
			xUnitFileSource := "xunit-reports/report.xml"
			xUnitContent := trimmedFileContentOrFatal(t, filepath.Join(sourceArtifactsBaseDir, xUnitFileSource))
			destinationArtifactsBaseDir := filepath.Join(subrepoDir, pipelinectxt.ArtifactsPath)
			checkFileContent(t, destinationArtifactsBaseDir, xUnitFileSource, xUnitContent)
			checkFilesExist(t, destinationArtifactsBaseDir, pipelinectxt.ArtifactsManifestFilename)

			bitbucketClient := tasktesting.BitbucketClientOrFatal(t, k8sClient, namespaceConfig.Name, *privateCertFlag)
			checkBuildStatus(t, bitbucketClient, odsContext.GitCommitSHA, bitbucket.BuildStatusInProgress)
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func TestStartTaskFailsWithoutSuccessfulPipelineRunOfSubrepo(t *testing.T) {
	k8sClient := newK8sClient(t)
	if err := runStartTask(
		ott.WithSourceWorkspace(
			t,
			"../testdata/workspaces/hello-world-app",
			func(c *ttr.WorkspaceConfig) error {
				_ = setupBitbucketRepoWithSubdirOrFatal(t, c, k8sClient)
				return nil
			},
		),
		func(c *ttr.TaskRunConfig) error {
			c.Params = append(c.Params, ttr.TektonParamsFromStringParams(map[string]string{
				"url":               bitbucketURLForWorkspace(c.WorkspaceConfigs["source"]),
				"git-full-ref":      "refs/heads/master",
				"project":           tasktesting.BitbucketProjectKey,
				"pipeline-run-name": "foo",
				"artifact-source":   "empty-repo",
			})...)
			return nil
		},
		ttr.ExpectFailure(),
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			want := "Pipeline runs with subrepos require a successful pipeline run artifact " +
				"for all checked out subrepo commits, however no such artifact was found"

			if !strings.Contains(logs.String(), want) {
				t.Fatalf("Want:\n%s\n\nGot:\n%s", want, logs.String())
			}
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func TestStartTaskClonesUsingLFS(t *testing.T) {
	var lfsFilename string
	var lfsFileHash [32]byte
	k8sClient := newK8sClient(t)
	if err := runStartTask(
		ott.WithSourceWorkspace(
			t,
			"../testdata/workspaces/hello-world-app",
			func(c *ttr.WorkspaceConfig) error {
				odsContext := tasktesting.SetupBitbucketRepo(
					t, k8sClient, namespaceConfig.Name, c.Dir, tasktesting.BitbucketProjectKey, *privateCertFlag,
				)
				tasktesting.EnableLfsOnBitbucketRepoOrFatal(t, filepath.Base(c.Dir), tasktesting.BitbucketProjectKey)
				lfsFilename = "lfspicture.jpg"
				lfsFileHash = tasktesting.UpdateBitbucketRepoWithLfsOrFatal(t, odsContext, c.Dir, tasktesting.BitbucketProjectKey, lfsFilename)
				return nil
			},
		),
		func(c *ttr.TaskRunConfig) error {
			c.Params = append(c.Params, ttr.TektonParamsFromStringParams(map[string]string{
				"url":               bitbucketURLForWorkspace(c.WorkspaceConfigs["source"]),
				"git-full-ref":      "refs/heads/master",
				"project":           tasktesting.BitbucketProjectKey,
				"pipeline-run-name": "foo",
			})...)
			return nil
		},
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			wsDir, odsContext := ott.GetSourceWorkspaceContext(t, config)
			checkODSContext(t, wsDir, odsContext)
			checkFileHash(t, wsDir, lfsFilename, lfsFileHash)
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func setupBitbucketRepoWithSubdirOrFatal(t *testing.T, c *ttr.WorkspaceConfig, k8sClient kubernetes.Interface) *pipelinectxt.ODSContext {
	// Setup sub-component
	tempDir, err := directory.CopyToTempDir(
		filepath.Join(projectpath.Root, "test", testdataWorkspacesPath, "hello-world-app"),
		c.Dir,
		"subcomponent-",
	)
	if err != nil {
		t.Fatal(err)
	}
	subCtxt := tasktesting.SetupBitbucketRepo(
		t, k8sClient, namespaceConfig.Name, tempDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
	)
	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	err = createStartODSYMLWithSubrepo(c.Dir, filepath.Base(tempDir))
	if err != nil {
		t.Fatal(err)
	}
	_ = tasktesting.SetupBitbucketRepo(
		t, k8sClient, namespaceConfig.Name, c.Dir, tasktesting.BitbucketProjectKey, *privateCertFlag,
	)
	return subCtxt
}

func bitbucketURLForWorkspace(c *ttr.WorkspaceConfig) string {
	bbURL := "http://ods-test-bitbucket-server.kind:7990"
	repoName := filepath.Base(c.Dir)
	return fmt.Sprintf("%s/scm/%s/%s.git", bbURL, tasktesting.BitbucketProjectKey, repoName)
}

func createStartODSYMLWithSubrepo(wsDir, repo string) error {
	o := &config.ODS{Repositories: []config.Repository{{Name: repo}}}
	return createODSYML(wsDir, o)
}

func createODSYML(wsDir string, o *config.ODS) error {
	y, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	filename := filepath.Join(wsDir, "ods.yaml")
	return os.WriteFile(filename, y, 0644)
}

func checkFileHash(t *testing.T, wsDir string, filename string, hash [32]byte) {
	filepath := filepath.Join(wsDir, filename)
	filecontent, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Want %s, but got nothing", filename)
	}
	filehash := sha256.Sum256(filecontent)
	if filehash != hash {
		t.Fatalf("Want %x, but got %x", hash, filehash)
	}
}

func checkODSContext(t *testing.T, repoDir string, want *pipelinectxt.ODSContext) {
	checkODSFileContent(t, repoDir, "component", want.Component)
	checkODSFileContent(t, repoDir, "git-commit-sha", want.GitCommitSHA)
	checkODSFileContent(t, repoDir, "git-full-ref", want.GitFullRef)
	checkODSFileContent(t, repoDir, "git-ref", want.GitRef)
	checkODSFileContent(t, repoDir, "git-url", want.GitURL)
	checkODSFileContent(t, repoDir, "namespace", want.Namespace)
	checkODSFileContent(t, repoDir, "pr-base", want.PullRequestBase)
	checkODSFileContent(t, repoDir, "pr-key", want.PullRequestKey)
	checkODSFileContent(t, repoDir, "project", want.Project)
	checkODSFileContent(t, repoDir, "repository", want.Repository)
}

func checkODSFileContent(t *testing.T, wsDir, filename, want string) {
	checkFileContent(t, filepath.Join(wsDir, pipelinectxt.BaseDir), filename, want)
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

func checkFilesExist(t *testing.T, wsDir string, wantFiles ...string) {
	for _, wf := range wantFiles {
		filename := filepath.Join(wsDir, wf)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Fatalf("Want %s, but got nothing", filename)
		}
	}
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func trimmedFileContentOrFatal(t *testing.T, filename string) string {
	c, err := getTrimmedFileContent(filename)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
