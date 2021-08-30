package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

type PipelineRunArtifact struct {
	// Name is the pipeline run name.
	Name string `json:"name"`
	// AggregateTaskStatus is the aggregate Tekton task status.
	AggregateTaskStatus string `json:"aggregateTaskStatus"`
}

type options struct {
	bitbucketAccessToken     string
	bitbucketURL             string
	consoleURL               string
	pipelineRunName          string
	aggregateTasksStatus     string
	nexusURL                 string
	nexusUsername            string
	nexusPassword            string
	nexusTemporaryRepository string
	nexusPermanentRepository string
	debug                    bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.bitbucketAccessToken, "bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	flag.StringVar(&opts.bitbucketURL, "bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	flag.StringVar(&opts.consoleURL, "console-url", os.Getenv("CONSOLE_URL"), "web console URL")
	flag.StringVar(&opts.pipelineRunName, "pipeline-run-name", "", "name of pipeline run")
	// See https://tekton.dev/docs/pipelines/pipelines/#using-aggregate-execution-status-of-all-tasks.
	// Possible values are: Succeeded, Failed, Completed, None.
	flag.StringVar(&opts.aggregateTasksStatus, "aggregate-tasks-status", "None", "aggregate status of all the tasks")
	flag.StringVar(&opts.nexusURL, "nexus-url", os.Getenv("NEXUS_URL"), "Nexus URL")
	flag.StringVar(&opts.nexusUsername, "nexus-username", os.Getenv("NEXUS_USERNAME"), "Nexus username")
	flag.StringVar(&opts.nexusPassword, "nexus-password", os.Getenv("NEXUS_PASSWORD"), "Nexus password")
	flag.StringVar(&opts.nexusTemporaryRepository, "nexus-temporary-repository", os.Getenv("NEXUS_TEMPORARY_REPOSITORY"), "Nexus temporary repository")
	flag.StringVar(&opts.nexusPermanentRepository, "nexus-permanent-repository", os.Getenv("NEXUS_PERMANENT_REPOSITORY"), "Nexus permanent repository")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "debug mode")
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

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	}

	fmt.Println("Setting Bitbucket build status ...")
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: opts.bitbucketAccessToken,
		BaseURL:  opts.bitbucketURL,
		Logger:   logger,
	})
	pipelineRunURL := fmt.Sprintf(
		"%s/k8s/ns/%s/tekton.dev~v1beta1~PipelineRun/%s/",
		opts.consoleURL,
		ctxt.Namespace,
		opts.pipelineRunName,
	)
	err = bitbucketClient.BuildStatusCreate(ctxt.GitCommitSHA, bitbucket.BuildStatusCreatePayload{
		State:       getBitbucketBuildStatus(opts.aggregateTasksStatus),
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
		BaseURL:    opts.nexusURL,
		Username:   opts.nexusUsername,
		Password:   opts.nexusPassword,
		Repository: opts.nexusTemporaryRepository,
		Logger:     logger,
	})
	if err != nil {
		log.Fatal(err)
	}

	if tasksSuccessful(opts.aggregateTasksStatus) {
		fmt.Println("Creating artifact of pipeline run ...")
		err := createPipelineRunArtifact(opts.pipelineRunName, opts.aggregateTasksStatus)
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
	return pipelinectxt.WriteJsonArtifact(pra, pipelinectxt.PipelineRunsPath, pra.Name+".json")
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
