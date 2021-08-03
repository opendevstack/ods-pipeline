package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

func main() {
	bitbucketAccessTokenFlag := flag.String("bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	bitbucketURLFlag := flag.String("bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	consoleURLFlag := flag.String("console-url", os.Getenv("CONSOLE_URL"), "web console URL")
	pipelineRunNameFlag := flag.String("pipeline-run-name", "", "name of pipeline run")
	// See https://tekton.dev/docs/pipelines/pipelines/#using-aggregate-execution-status-of-all-tasks.
	// Possible values are: Succeeded, Failed, Completed, None.
	aggregateTasksStatusFlag := flag.String("aggregate-tasks-status", "None", "aggregate status of all the tasks")
	nexusURLFlag := flag.String("nexus-url", os.Getenv("NEXUS_URL"), "Nexus URL")
	nexusUsernameFlag := flag.String("nexus-username", os.Getenv("NEXUS_USERNAME"), "Nexus username")
	nexusPasswordFlag := flag.String("nexus-password", os.Getenv("NEXUS_PASSWORD"), "Nexus password")
	nexusRepositoryFlag := flag.String("nexus-repository", "ods-pipelines", "Nexus repository")
	flag.Parse()

	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(".")
	if err != nil {
		log.Fatalf(
			"Unable to continue as pipeline context cannot be read: %s.\n"+
				"Bitbucket build status will not be set and no artifacts will be uploaded to Nexus.",
			err,
		)
	}

	// Set Bitbucket build status
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: *bitbucketAccessTokenFlag,
		BaseURL:  *bitbucketURLFlag,
	})
	pipelineRunURL := fmt.Sprintf(
		"%s/k8s/ns/%s/tekton.dev~v1beta1~PipelineRun/%s/",
		*consoleURLFlag,
		ctxt.Namespace,
		*pipelineRunNameFlag,
	)
	err = bitbucketClient.BuildStatusCreate(ctxt.GitCommitSHA, bitbucket.BuildStatusCreatePayload{
		State:       getBuildStatus(*aggregateTasksStatusFlag),
		Key:         ctxt.GitCommitSHA,
		Name:        ctxt.GitCommitSHA,
		URL:         pipelineRunURL,
		Description: "ODS Pipeline Build",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Upload all files in .ods/artifacts to Nexus in a folder named like the Git commit SHA.
	nexusClient, err := nexus.NewClient(
		*nexusURLFlag,
		*nexusUsernameFlag,
		*nexusPasswordFlag,
		*nexusRepositoryFlag,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Read files from .ods/artifacts
	artifactsMap, err := pipelinectxt.ReadArtifactsDir()
	if err != nil {
		log.Fatal(err)
	}

	for artifactsSubDir, files := range artifactsMap {
		for _, filename := range files {
			err = nexusClient.Upload(fmt.Sprintf("/%s/%s/%s/%s", ctxt.Project, ctxt.Repository, ctxt.GitCommitSHA, artifactsSubDir), ".ods/artifacts/"+artifactsSubDir+"/"+filename)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func getBuildStatus(aggregateTasksStatus string) string {
	// Meaning of aggregateTasksStatus values:
	// Succeeded: all tasks have succeeded.
	// Failed: one ore more tasks failed.
	// Completed: all tasks completed successfully including one or more skipped tasks.
	// None: no aggregate execution status available (i.e. none of the above),
	// one or more tasks could be pending/running/cancelled/timedout.
	// See https://tekton.dev/docs/pipelines/pipelines/#using-aggregate-execution-status-of-all-tasks.
	if aggregateTasksStatus == "Succeeded" || aggregateTasksStatus == "Completed" {
		return "SUCCESSFUL"
	} else {
		return "FAILED"
	}
}
