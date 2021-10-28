package tasks

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
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
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey,
					)
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
			"set bitbucket build status to successful and upload artifacts to temporary Nexus repository": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app-with-artifacts"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey,
					)
					// Pretend there is alredy a coverage report in Nexus
					// this assures the safeguard is working to avoid duplicate upload.
					// TODO: assure the safeguard is actually invoked by checking the logs.
					t.Log("Uploading coverage artifact to Nexus and writing manifest")
					nexusClient := tasktesting.NexusClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					err := nexusClient.Upload(
						nexus.TemporaryRepositoryDefault,
						nexus.ArtifactGroup(ctxt.ODS, pipelinectxt.CodeCoveragesDir),
						filepath.Join(wsDir, pipelinectxt.CodeCoveragesPath, "coverage.out"),
					)
					if err != nil {
						t.Fatal(err)
					}
					am := pipelinectxt.ArtifactsManifest{
						SourceRepository: nexus.TemporaryRepositoryDefault,
						Artifacts: []pipelinectxt.ArtifactInfo{
							{
								Directory: pipelinectxt.CodeCoveragesDir,
								Name:      "coverage.out",
							},
						},
					}
					err = pipelinectxt.WriteJsonArtifact(
						am,
						filepath.Join(wsDir, pipelinectxt.ArtifactsPath),
						pipelinectxt.ArtifactsManifestFilename,
					)
					if err != nil {
						t.Fatal(err)
					}

					ctxt.Params = map[string]string{
						"pipeline-run-name":      "foo",
						"aggregate-tasks-status": "Succeeded",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusSuccessful)
					checkArtifactsAreInNexus(t, ctxt, nexus.TemporaryRepositoryDefault)
				},
			},
			"set bitbucket build status to successful and upload artifacts to permanent Nexus repository": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app-with-artifacts"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupBitbucketRepo(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, wsDir, tasktesting.BitbucketProjectKey,
					)
					err := createFinishODSYML(wsDir)
					if err != nil {
						t.Fatal(err)
					}
					ctxt.Params = map[string]string{
						"pipeline-run-name":      "foo",
						"aggregate-tasks-status": "Succeeded",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					bitbucketClient := tasktesting.BitbucketClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					checkBuildStatus(t, bitbucketClient, ctxt.ODS.GitCommitSHA, bitbucket.BuildStatusSuccessful)
					checkArtifactsAreInNexus(t, ctxt, nexus.PermanentRepositoryDefault)
				},
			},
		},
	)
}

func createFinishODSYML(wsDir string) error {
	o := &config.ODS{
		Environments: []config.Environment{
			{
				Name:  "dev",  // use "dev" because that is set in the context, but ...
				Stage: "prod", // set the stage to "prod" to simulate a production pipeline.
			},
		},
	}
	return createODSYML(wsDir, o)
}

func checkArtifactsAreInNexus(t *testing.T, ctxt *tasktesting.TaskRunContext, targetRepository string) {

	nexusClient := tasktesting.NexusClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)

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
		group := nexus.ArtifactGroup(ctxt.ODS, artifactsSubDir)

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
