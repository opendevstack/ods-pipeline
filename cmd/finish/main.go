package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

func main() {
	bitbucketAccessTokenFlag := flag.String("bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	bitbucketURLFlag := flag.String("bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	consoleURLFlag := flag.String("console-url", "", "web console URL")
	pipelineRunNameFlag := flag.String("pipeline-run-name", "", "name of pipeline run")
	aggregateTasksStatusFlag := flag.String("aggregate-tasks-status", "None", "aggregate status of all the tasks")
	flag.Parse()

	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(".")
	if err != nil {
		panic(err.Error())
	}

	// Set Bitbucket build status
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		Timeout:    10 * time.Second,
		APIToken:   *bitbucketAccessTokenFlag,
		MaxRetries: 2,
		BaseURL:    *bitbucketURLFlag,
	})
	pipelineRunURL := fmt.Sprintf(
		"%s/k8s/ns/%s/tekton.dev~v1beta1~PipelineRun/%s/",
		*consoleURLFlag,
		ctxt.Namespace,
		*pipelineRunNameFlag,
	)
	err = bitbucketClient.BuildStatusPost(ctxt.GitCommitSHA, bitbucket.BuildStatusPostPayload{
		State:       getBuildStatus(aggregateTasksStatusFlag),
		Key:         ctxt.GitCommitSHA,
		Name:        ctxt.GitCommitSHA,
		URL:         pipelineRunURL,
		Description: "ODS Pipeline Build",
	})
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Upload all files in .ods/artifacts to Nexus in a folder named like the Git commit SHA.

}

func getBuildStatus(aggregateTasksStatusFlag *string) string {

	if *aggregateTasksStatusFlag == "Succeeded" {
		return "SUCCESSFUL"
	} else {
		return "FAILED"
	}
}
