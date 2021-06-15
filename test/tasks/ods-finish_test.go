package tasks

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

// TODO:
// Read the from configmap and secret files
// from test/testdata/deploy/cd-kind
const (
	nexusURL        = "http://localhost:8081"
	nexusUser       = "developer"
	nexusPassword   = "s3cr3t"
	nexusRepository = "ods-pipelines"
)

func TestTaskODSFinish(t *testing.T) {

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
		"set bitbucket build status to successful and artifacts are in Nexus": {
			WorkspaceDirMapping: map[string]string{"source": "hello-world-app-with-artifacts"},
			PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				ctxt.ODS = tasktesting.SetupBitbucketRepo(t, c.KubernetesClientSet, ns, wsDir, bitbucketProjectKey)
				ctxt.Params = map[string]string{
					"image":                  "localhost:5000/ods/finish:latest",
					"console-url":            "http://example.com",
					"pipeline-run-name":      "foo",
					"aggregate-tasks-status": "Succeeded",
				}
			},
			WantRunSuccess: true,
			PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {

				// TODO: Check Bitbucket build status is successful

				checkArtifactsAreInNexus(t, ctxt)
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             "ClusterTask",       // could be read from task definition
				TaskName:                "ods-finish-v0-1-0", // could be read from task definition
				Clients:                 c,
				Namespace:               ns,
				Timeout:                 5 * time.Minute, // depending on  the task we may need to increase or decrease it
				AlwaysKeepTmpWorkspaces: *alwaysKeepTmpWorkspacesFlag,
			})

		})

	}
}

func checkArtifactsAreInNexus(t *testing.T, ctxt *tasktesting.TaskRunContext) {

	nexusClient, err := nexus.NewClient(
		nexusURL,
		nexusUser,
		nexusPassword,
		nexusRepository,
	)
	if err != nil {
		t.Fatal(err)
	}

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
			url := fmt.Sprintf("%s/repository/%s/%s/%s/%s/%s/%s", nexusURL, nexusRepository, ctxt.ODS.Project, ctxt.ODS.Repository, ctxt.ODS.GitCommitSHA, artifactsSubDir, file)

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
		artifactURLs, err := nexusClient.URLs(group)
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
