package e2e

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/ods-pipeline/internal/tasktesting"
	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	"github.com/opendevstack/ods-pipeline/pkg/nexus"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"k8s.io/client-go/kubernetes"

	ott "github.com/opendevstack/ods-pipeline/pkg/odstasktest"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
)

func runFinishTask(opts ...ttr.TaskRunOpt) error {
	return ttr.RunTask(append([]ttr.TaskRunOpt{
		ttr.InNamespace(namespaceConfig.Name),
		ttr.UsingTask("ods-pipeline-finish"),
	}, opts...)...)
}

func TestFinishTaskSetsBitbucketStatusToFailed(t *testing.T) {
	k8sClient := newK8sClient(t)
	if err := runFinishTask(
		withBitbucketSourceWorkspace(t, "../testdata/workspaces/hello-world-app-with-artifacts", k8sClient, namespaceConfig.Name),
		ttr.WithStringParams(map[string]string{
			"pipeline-run-name":      "foo",
			"aggregate-tasks-status": "None",
		}),
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			_, odsContext := ott.GetSourceWorkspaceContext(t, config)
			bitbucketClient := tasktesting.BitbucketClientOrFatal(t, k8sClient, namespaceConfig.Name, *privateCertFlag)
			checkBuildStatus(t, bitbucketClient, odsContext.GitCommitSHA, bitbucket.BuildStatusFailed)
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func TestFinishTaskSetsBitbucketStatusToSuccessfulAndUploadsArtifactsToNexus(t *testing.T) {
	k8sClient := newK8sClient(t)
	if err := runFinishTask(
		ott.WithSourceWorkspace(
			t,
			"../testdata/workspaces/hello-world-app-with-artifacts",
			func(c *ttr.WorkspaceConfig) error {
				odsContext := tasktesting.SetupBitbucketRepo(
					t, k8sClient, namespaceConfig.Name, c.Dir, tasktesting.BitbucketProjectKey, *privateCertFlag,
				)
				// Pretend there is alredy a coverage report in Nexus.
				// This assures the safeguard is working to avoid duplicate upload.
				t.Log("Uploading coverage artifact to Nexus and writing manifest")
				nexusClient := tasktesting.NexusClientOrFatal(t, k8sClient, namespaceConfig.Name, *privateCertFlag)
				if _, err := nexusClient.Upload(
					nexus.TestTemporaryRepository,
					pipelinectxt.ArtifactGroup(odsContext, pipelinectxt.CodeCoveragesDir),
					filepath.Join(c.Dir, pipelinectxt.CodeCoveragesPath, "coverage.out"),
				); err != nil {
					t.Fatal(err)
				}
				am := pipelinectxt.NewArtifactsManifest(
					nexus.TestTemporaryRepository,
					pipelinectxt.ArtifactInfo{
						Directory: pipelinectxt.CodeCoveragesDir,
						Name:      "coverage.out",
					},
				)
				if err := pipelinectxt.WriteJsonArtifact(
					am,
					filepath.Join(c.Dir, pipelinectxt.ArtifactsPath),
					pipelinectxt.ArtifactsManifestFilename,
				); err != nil {
					t.Fatal(err)
				}
				return nil
			},
		),
		ttr.WithStringParams(map[string]string{
			"pipeline-run-name":      "foo",
			"aggregate-tasks-status": "Succeeded",
			"artifact-target":        nexus.TestTemporaryRepository,
		}),
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			_, odsContext := ott.GetSourceWorkspaceContext(t, config)
			bitbucketClient := tasktesting.BitbucketClientOrFatal(t, k8sClient, namespaceConfig.Name, *privateCertFlag)
			checkBuildStatus(t, bitbucketClient, odsContext.GitCommitSHA, bitbucket.BuildStatusSuccessful)
			checkArtifactsAreInNexus(t, k8sClient, odsContext, nexus.TestTemporaryRepository)

			wantLogMsg := "Artifact \"coverage.out\" is already present in Nexus repository"
			if !strings.Contains(logs.String(), wantLogMsg) {
				t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, logs.String())
			}
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func TestFinishTaskStopsGracefullyWhenContextCannotBeRead(t *testing.T) {
	if err := runFinishTask(
		ott.WithSourceWorkspace(t, "../testdata/workspaces/empty"),
		ttr.WithStringParams(map[string]string{
			"pipeline-run-name":      "foo",
			"aggregate-tasks-status": "None",
		}),
		ttr.ExpectFailure(),
		ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun, logs bytes.Buffer) {
			want := "Unable to continue as pipeline context cannot be read"
			if !strings.Contains(logs.String(), want) {
				t.Fatalf("Want:\n%s\n\nGot:\n%s", want, logs.String())
			}
		}),
	); err != nil {
		t.Fatal(err)
	}
}

func checkArtifactsAreInNexus(t *testing.T, k8sClient kubernetes.Interface, odsContext *pipelinectxt.ODSContext, targetRepository string) {

	nexusClient := tasktesting.NexusClientOrFatal(t, k8sClient, namespaceConfig.Name, *privateCertFlag)

	// List of expected artifacts to have been uploaded to Nexus
	artifactsMap := map[string][]string{
		pipelinectxt.XUnitReportsDir: {"report.xml"},
		// exclude coverage as we pretend it has been uploaded earlier already
		// pipelinectxt.CodeCoveragesDir: {"coverage.out"},
		pipelinectxt.SonarAnalysisDir: {"analysis-report.md", "issues-report.csv"},
	}

	for artifactsSubDir, files := range artifactsMap {

		filesCountInSubDir := len(artifactsMap[artifactsSubDir])

		// e.g: "/ODSPIPELINETEST/workspace-190880007/935e5229b084dd60d44a5eddd2d023720ec153c1/xunit-reports"
		group := pipelinectxt.ArtifactGroup(odsContext, artifactsSubDir)

		// The test is so fast that, when we reach this line, the artifacts could still being uploaded to Nexus
		artifactURLs := waitForArtifacts(t, nexusClient, targetRepository, group, filesCountInSubDir, 5*time.Second)
		if len(artifactURLs) != filesCountInSubDir {
			t.Fatalf("Got: %d artifacts in subdir %s, want: %d.", len(artifactURLs), artifactsMap[artifactsSubDir], filesCountInSubDir)
		}

		for _, file := range files {

			// e.g. "http://localhost:8081/repository/ods-pipelines/ODSPIPELINETEST/workspace-866704509/b1415e831b4f5b24612abf24499663ddbff6babb/xunit-reports/report.xml"
			// note that the "group" value already has a leading slash!
			url := fmt.Sprintf("%s/repository/%s%s/%s", nexusClient.URL(), targetRepository, group, file)

			if !contains(artifactURLs, url) {
				t.Fatalf("Artifact %s with URL %+v not found in Nexus under any of the following URLs: %v", file, url, artifactURLs)
			}
		}

	}
}

func waitForArtifacts(t *testing.T, nexusClient *nexus.Client, targetRepository, group string, expectedArtifactsCount int, timeout time.Duration) []string {
	start := time.Now().UTC()
	elapsed := time.Since(start)
	artifactURLs := []string{}

	for elapsed < timeout {
		artifactURLs, err := nexusClient.Search(targetRepository, group)
		if err != nil {
			t.Fatal(err)
		}

		if len(artifactURLs) == expectedArtifactsCount {
			return artifactURLs
		}

		log.Printf("Artifacts are not yet available in Nexus...\n")
		time.Sleep(1 * time.Second)

		elapsed = time.Since(start)
	}

	log.Printf("Time out reached.\n")
	return artifactURLs
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
