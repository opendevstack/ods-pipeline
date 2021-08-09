package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

type PipelineRunArtifact struct {
	// Name is the pipeline run name.
	Name string `json:"name"`
	// AggregateTaskStatus is the aggregate Tekton task status.
	AggregateTaskStatus string `json:"aggregateTaskStatus"`
}

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
	nexusTemporaryRepositoryFlag := flag.String("nexus-temporary-repository", os.Getenv("NEXUS_TEMPORARY_REPOSITORY"), "Nexus temporary repository")
	//nexusPermanentRepositoryFlag := flag.String("nexus-permanent-repository", os.Getenv("NEXUS_PERMANENT_REPOSITORY"), "Nexus permanent repository")
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

	fmt.Println("Setting Bitbucket build status ...")
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
		State:       getBitbucketBuildStatus(*aggregateTasksStatusFlag),
		Key:         ctxt.GitCommitSHA,
		Name:        ctxt.GitCommitSHA,
		URL:         pipelineRunURL,
		Description: "ODS Pipeline Build",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Handling artifacts ...")
	nexusClient, err := nexus.NewClient(&nexus.ClientConfig{
		BaseURL:    *nexusURLFlag,
		Username:   *nexusUsernameFlag,
		Password:   *nexusPasswordFlag,
		Repository: *nexusTemporaryRepositoryFlag,
	})
	if err != nil {
		log.Fatal(err)
	}

	if tasksSuccessful(*aggregateTasksStatusFlag) {
		fmt.Println("Creating artifact of pipeline run ...")
		err := createPipelineRunArtifact(*pipelineRunNameFlag, *aggregateTasksStatusFlag)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Uploading artifacts to Nexus ...")
		artifactsMap, err := pipelinectxt.ReadArtifactsDir()
		if err != nil {
			log.Fatal(err)
		}

		for artifactsSubDir, files := range artifactsMap {
			for _, filename := range files {
				nexusGroup := fmt.Sprintf("/%s/%s/%s/%s", ctxt.Project, ctxt.Repository, ctxt.GitCommitSHA, artifactsSubDir)
				localFile := filepath.Join(pipelinectxt.ArtifactsPath, artifactsSubDir, filename)
				fmt.Printf("Uploading %s to Nexus group %s ...\n", localFile, nexusGroup)
				err = nexusClient.Upload(nexusGroup, localFile)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		log.Println("No artifacts are uploaded to Nexus as one or more tasks failed.")
	}

}

func createPipelineRunArtifact(pipelineRunName, aggregateTasksStatus string) error {
	pra := PipelineRunArtifact{
		Name:                pipelineRunName,
		AggregateTaskStatus: aggregateTasksStatus,
	}
	j, err := json.Marshal(pra)
	if err != nil {
		return fmt.Errorf("could not marshal pipeline run artifact: %w", err)
	}
	err = os.MkdirAll(pipelinectxt.PipelineRunsPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create pipeline run artifact directory: %w", err)
	}
	filename := filepath.Join(pipelinectxt.PipelineRunsPath, pra.Name+".json")
	err = ioutil.WriteFile(filename, j, 0644)
	if err != nil {
		return fmt.Errorf("could not write pipeline run artifact: %w", err)
	}
	return nil
}

// getBitbucketBuildStatus returns a build status for use with Bitbucket based
// on the aggregate Tekton tasks status.
// See https://developer.atlassian.com/server/bitbucket/how-tos/updating-build-status-for-commits/.
func getBitbucketBuildStatus(aggregateTasksStatus string) string {
	if tasksSuccessful(aggregateTasksStatus) {
		return bitbucket.BuildStatusSuccessful
	} else {
		return bitbucket.BuildStatusFailed
	}
}

// tasksSuccessful returns true if no task failed.
func tasksSuccessful(aggregateTasksStatus string) bool {
	// Meaning of aggregateTasksStatus values:
	// Succeeded: all tasks have succeeded.
	// Failed: one ore more tasks failed.
	// Completed: all tasks completed successfully including one or more skipped tasks.
	// None: no aggregate execution status available (i.e. none of the above),
	// one or more tasks could be pending/running/cancelled/timedout.
	// See https://tekton.dev/docs/pipelines/pipelines/#using-aggregate-execution-status-of-all-tasks.
	return aggregateTasksStatus == "Succeeded" || aggregateTasksStatus == "Completed"
}
