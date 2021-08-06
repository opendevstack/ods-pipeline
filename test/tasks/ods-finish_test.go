package tasks

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSFinish(t *testing.T) {
	runTaskTestCases(t,
		"ods-finish",
		map[string]tasktesting.TestCase{
			"set bitbucket build status to failed": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app-with-artifacts"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupBitbucketRepo(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, bitbucketProjectKey)
					ctxt.Params = map[string]string{
						"pipeline-run-name":      "foo",
						"aggregate-tasks-status": "None",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusFailed)
				},
			},
			"set bitbucket build status to successful and artifacts are in Nexus": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app-with-artifacts"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupBitbucketRepo(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, bitbucketProjectKey)
					ctxt.Params = map[string]string{
						"pipeline-run-name":      "foo",
						"aggregate-tasks-status": "Succeeded",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusSuccessful)
					checkArtifactsAreInNexus(t, ctxt)
				},
			},
		},
	)
}

func checkArtifactsAreInNexus(t *testing.T, ctxt *tasktesting.TaskRunContext) {

	nexusClient := tasktesting.NexusClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)

	// List of expected artifacts to have been uploaded to Nexus
	artifactsMap := map[string][]string{
		"xunit-reports":      {"report.xml"},
		"code-coverage":      {"coverage.out"},
		"sonarqube-analysis": {"analysis-report.md", "issues-report.csv"},
	}

	for artifactsSubDir, files := range artifactsMap {

		filesCountInSubDir := len(artifactsMap[artifactsSubDir])

		// e.g: "/ODSPIPELINETEST/workspace-190880007/935e5229b084dd60d44a5eddd2d023720ec153c1/xunit-reports"
		group := fmt.Sprintf("/%s/%s/%s/%s", ctxt.ODS.Project, ctxt.ODS.Repository, ctxt.ODS.GitCommitSHA, artifactsSubDir)

		// The test is so fast that, when we reach this line, the artifacts could still being uploaded to Nexus
		artifactURLs := waitForArtifacts(t, nexusClient, group, filesCountInSubDir, 5*time.Second)
		if len(artifactURLs) != filesCountInSubDir {
			t.Fatalf("Got: %d artifacts in subdir %s, want: %d.", len(artifactURLs), artifactsMap[artifactsSubDir], filesCountInSubDir)
		}

		for _, file := range files {

			// e.g. "http://localhost:8081/repository/ods-pipelines/ODSPIPELINETEST/workspace-866704509/b1415e831b4f5b24612abf24499663ddbff6babb/xunit-reports/report.xml"
			url := fmt.Sprintf("%s/repository/%s/%s/%s/%s/%s/%s", nexusClient.URL(), nexusClient.Repository(), ctxt.ODS.Project, ctxt.ODS.Repository, ctxt.ODS.GitCommitSHA, artifactsSubDir, file)

			if !contains(artifactURLs, url) {
				t.Fatalf("Artifact %s with URL %+v not found in Nexus under any of the following URLs: %v", file, url, artifactURLs)
			}
		}

	}
}

func waitForArtifacts(t *testing.T, nexusClient *nexus.Client, group string, expectedArtifactsCount int, timeout time.Duration) []string {

	start := time.Now().UTC()
	elapsed := time.Since(start)
	artifactURLs := []string{}

	for elapsed < timeout {
		artifactURLs, err := nexusClient.Search(group)
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
