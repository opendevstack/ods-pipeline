package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/notification"
	"github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
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
	checkoutDir := "."

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
	flag.StringVar(&opts.nexusTemporaryRepository, nexus.TemporaryRepositoryDefault, os.Getenv("NEXUS_TEMPORARY_REPOSITORY"), "Nexus temporary repository")
	flag.StringVar(&opts.nexusPermanentRepository, nexus.PermanentRepositoryDefault, os.Getenv("NEXUS_PERMANENT_REPOSITORY"), "Nexus permanent repository")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "debug mode")
	flag.Parse()

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(checkoutDir)
	if err != nil {
		log.Fatalf(
			"Unable to continue as pipeline context cannot be read: %s.\n"+
				"Bitbucket build status will not be set and no artifacts will be uploaded to Nexus.",
			err,
		)
	}

	logger.Infof("Setting Bitbucket build status ...")
	bitbucketClient, err := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: opts.bitbucketAccessToken,
		BaseURL:  opts.bitbucketURL,
		Logger:   logger,
	})
	if err != nil {
		log.Fatal("bitbucket client:", err)
	}

	prURL, err := tekton.PipelineRunURL(opts.consoleURL, ctxt.Namespace, opts.pipelineRunName)
	if err != nil {
		log.Fatal("pipeline run URL:", err)
	}

	err = bitbucketClient.BuildStatusCreate(ctxt.GitCommitSHA, bitbucket.BuildStatusCreatePayload{
		State:       getBitbucketBuildStatus(opts.aggregateTasksStatus),
		Key:         ctxt.GitCommitSHA,
		Name:        ctxt.GitCommitSHA,
		URL:         prURL,
		Description: "ODS Pipeline Build",
	})
	if err != nil {
		log.Fatal(err)
	}

	nexusClient, err := nexus.NewClient(&nexus.ClientConfig{
		BaseURL:  opts.nexusURL,
		Username: opts.nexusUsername,
		Password: opts.nexusPassword,
		Logger:   logger,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = handleArtifacts(logger, nexusClient, opts, checkoutDir, ctxt)
	if err != nil {
		log.Fatal(err)
	}

	kubernetesClient, err := kubernetes.NewInClusterClient(&kubernetes.ClientConfig{
		Namespace: ctxt.Namespace,
	})
	if err != nil {
		log.Fatalf("couldn't create kubernetes client: %s", err)
	}

	ctx := context.TODO()
	notificationConfig, err := notification.ReadConfigFromConfigMap(ctx, kubernetesClient)
	if err != nil {
		log.Fatalf("Notification config could not be read: %s", err)
	}

	notificationClient, err := notification.NewClient(notification.ClientConfig{
		Namespace:          ctxt.Namespace,
		NotificationConfig: notificationConfig,
	})
	if err != nil {
		log.Fatal(err)
	}

	if notificationClient.ShouldNotify(opts.aggregateTasksStatus) {
		err = notificationClient.CallWebhook(ctx, notification.PipelineRunResult{
			PipelineRunName: opts.pipelineRunName,
			PipelineRunURL:  prURL,
			OverallStatus:   opts.aggregateTasksStatus,
			ODSContext:      ctxt,
		})
		if err != nil {
			log.Printf("Calling notification webhook failed: %s", err)
		}
	}
}

// handleArtifacts figures out what to do with the artifacts stored underneath
// the checkout dir. If the previous Tekton tasks succeeded, then the artifacts
// are uploaded. The target repository is determined by the environment stage.
func handleArtifacts(
	logger logging.LeveledLoggerInterface,
	nexusClient nexus.ClientInterface,
	opts options,
	checkoutDir string,
	ctxt *pipelinectxt.ODSContext) error {
	logger.Infof("Handling artifacts ...")

	logger.Infof("Creating artifact of pipeline run ...")
	err := createPipelineRunArtifact(checkoutDir, opts.pipelineRunName, opts.aggregateTasksStatus)
	if err != nil {
		return fmt.Errorf("cannot create pipeline run artifact: %w", err)
	}

	odsConfig, err := config.ReadFromDir(checkoutDir)
	if err != nil {
		return fmt.Errorf("cannot read ods config: %w", err)
	}

	logger.Infof("Uploading artifacts to Nexus ...")
	// Use temporary storage for DEV stage and permanenet repository for
	// QA and PROD stage environments.
	nexusRepository := opts.nexusTemporaryRepository
	if len(ctxt.Environment) > 0 {
		env, err := odsConfig.Environment(ctxt.Environment)
		if err != nil {
			return fmt.Errorf("cannot determine environment: %w", err)
		}
		logger.Debugf("Selected environment %s.", env.Name)
		if env.Stage != config.DevStage {
			logger.Debugf(
				"Setting target Nexus repository to %s as environment is not a dev stage environment.",
				opts.nexusPermanentRepository,
			)
			nexusRepository = opts.nexusPermanentRepository
		}
	}
	err = uploadArtifacts(logger, nexusClient, nexusRepository, checkoutDir, ctxt, opts)
	if err != nil {
		return fmt.Errorf("cannot upload artifacts of main repository: %w", err)
	}
	if len(odsConfig.Repositories) > 0 {
		for _, subrepo := range odsConfig.Repositories {
			subrepoCheckoutDir := filepath.Join(checkoutDir, pipelinectxt.SubreposPath, subrepo.Name)
			subrepoCtxt := &pipelinectxt.ODSContext{}
			err := subrepoCtxt.ReadCache(subrepoCheckoutDir)
			if err != nil {
				return fmt.Errorf("cannot read cache of subrepository %s: %w", subrepo.Name, err)
			}
			err = uploadArtifacts(logger, nexusClient, nexusRepository, subrepoCheckoutDir, subrepoCtxt, opts)
			if err != nil {
				return fmt.Errorf("cannot upload artifacts of subrepository %s: %w", subrepo.Name, err)
			}
		}
	}

	return nil
}

// uploadArtifacts uploads artifacts stored in checkoutDir to given nexusRepository.
func uploadArtifacts(
	logger logging.LeveledLoggerInterface,
	nexusClient nexus.ClientInterface,
	nexusRepository, checkoutDir string,
	ctxt *pipelinectxt.ODSContext,
	opts options) error {
	logger.Infof("Handling artifacts in %s ...\n", checkoutDir)
	artifactsDir := filepath.Join(checkoutDir, pipelinectxt.ArtifactsPath)
	artifactsMap, err := pipelinectxt.ReadArtifactsDir(artifactsDir)
	if err != nil {
		return err
	}
	am, err := pipelinectxt.ReadArtifactsManifestFromFile(
		filepath.Join(checkoutDir, pipelinectxt.ArtifactsPath, pipelinectxt.ArtifactsManifestFilename),
	)
	if err != nil {
		return err
	}
	for artifactsSubDir, files := range artifactsMap {
		for _, filename := range files {
			if am.Contains(nexusRepository, artifactsSubDir, filename) {
				logger.Infof("Artifact %s is already present in Nexus repository %s.", filename, nexusRepository)
			} else {
				nexusGroup := artifactGroup(ctxt, artifactsSubDir, opts)
				localFile := filepath.Join(checkoutDir, pipelinectxt.ArtifactsPath, artifactsSubDir, filename)
				logger.Infof("Uploading %s to Nexus repository %s, group %s ...", localFile, nexusRepository, nexusGroup)
				link, err := nexusClient.Upload(nexusRepository, nexusGroup, localFile)
				if err != nil {
					return err
				}
				logger.Infof("Successfully uploaded %s to %s", localFile, link)
			}
		}
	}
	return nil
}

func artifactGroup(ctxt *pipelinectxt.ODSContext, artifactsSubDir string, opts options) string {
	if !tasksSuccessful(opts.aggregateTasksStatus) {
		artifactsSubDir = fmt.Sprintf("failed-%s-artifacts/%s", opts.pipelineRunName, artifactsSubDir)
	}
	return pipelinectxt.ArtifactGroup(ctxt, artifactsSubDir)
}

func createPipelineRunArtifact(checkoutDir string, pipelineRunName, aggregateTasksStatus string) error {
	pra := PipelineRunArtifact{
		Name:                pipelineRunName,
		AggregateTaskStatus: aggregateTasksStatus,
	}
	writeDir := filepath.Join(checkoutDir, pipelinectxt.PipelineRunsPath)
	return pipelinectxt.WriteJsonArtifact(pra, writeDir, pra.Name+".json")
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
