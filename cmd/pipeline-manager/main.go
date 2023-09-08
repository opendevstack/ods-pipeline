package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	kubernetesClient "github.com/opendevstack/ods-pipeline/internal/kubernetes"
	"github.com/opendevstack/ods-pipeline/internal/manager"
	tektonClient "github.com/opendevstack/ods-pipeline/internal/tekton"
	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
)

const (
	namespaceFile            = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	namespaceSuffix          = "-cd"
	repoBaseEnvVar           = "REPO_BASE"
	tokenEnvVar              = "ACCESS_TOKEN"
	webhookSecretEnvVar      = "WEBHOOK_SECRET"
	storageProvisionerEnvVar = "ODS_STORAGE_PROVISIONER"
	storageClassNameEnvVar   = "ODS_STORAGE_CLASS_NAME"
	storageClassNameDefault  = "standard"
	storageSizeEnvVar        = "ODS_STORAGE_SIZE"
	storageSizeDefault       = "2Gi"
	pruneMinKeepHoursEnvVar  = "ODS_PRUNE_MIN_KEEP_HOURS"
	pruneMinKeepHoursDefault = 48
	pruneMaxKeepRunsEnvVar   = "ODS_PRUNE_MAX_KEEP_RUNS"
	pruneMaxKeepRunsDefault  = 20
	initialWatchWait         = 10 * time.Second
	// Allow a few concurrent pipeline triggers before blocking.
	channelBufferSize = 5
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	err := serve()
	if err != nil {
		log.Fatalln(err)
	}
}

func serve() error {
	var logger logging.LeveledLoggerInterface
	if os.Getenv("DEBUG") == "true" {
		logger = &logging.LeveledLogger{Timestamp: true, Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Timestamp: true, Level: logging.LevelInfo}
	}
	logger.Infof("Booting ...")

	repoBase := os.Getenv(repoBaseEnvVar)
	if repoBase == "" {
		return fmt.Errorf("%s must be set", repoBaseEnvVar)
	}

	token := os.Getenv(tokenEnvVar)
	if token == "" {
		return fmt.Errorf("%s must be set", tokenEnvVar)
	}

	webhookSecret := os.Getenv(webhookSecretEnvVar)
	if webhookSecret == "" {
		return fmt.Errorf("%s must be set", webhookSecretEnvVar)
	}

	storageProvisioner := readStringFromEnvVar(storageProvisionerEnvVar, "")

	storageClassName := readStringFromEnvVar(storageClassNameEnvVar, storageClassNameDefault)

	storageSize := readStringFromEnvVar(storageSizeEnvVar, storageSizeDefault)

	pruneMinKeepHours, err := readIntFromEnvVar(
		pruneMinKeepHoursEnvVar, pruneMinKeepHoursDefault,
	)
	if err != nil {
		return err
	}
	pruneMaxKeepRuns, err := readIntFromEnvVar(
		pruneMaxKeepRunsEnvVar, pruneMaxKeepRunsDefault,
	)
	if err != nil {
		return err
	}

	namespace, err := getFileContent(namespaceFile)
	if err != nil {
		return err
	}

	project := strings.TrimSuffix(namespace, namespaceSuffix)

	// Initialize Kubernetes client.
	kClient, err := kubernetesClient.NewInClusterClient(&kubernetesClient.ClientConfig{
		Namespace: namespace,
	})
	if err != nil {
		return err
	}

	// Initialize Tekton tClient.
	tClient, err := tektonClient.NewInClusterClient(&tektonClient.ClientConfig{
		Namespace: namespace,
	})
	if err != nil {
		return err
	}

	// Initialize Bitbucket client.
	bitbucketClient, err := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: token,
		BaseURL:  strings.TrimSuffix(repoBase, "/scm"),
	})
	if err != nil {
		return fmt.Errorf("bitbucket client: %w", err)
	}

	// triggeredReposChan is used to communicate repos for which pipelines
	// have been triggered between receiver and pruner.
	triggeredReposChan := make(chan string, channelBufferSize)
	// triggeredPipelinesChan is used to communicate triggered pipelines from
	// the receiver to the scheduler.
	triggeredPipelinesChan := make(chan manager.PipelineConfig, channelBufferSize)
	// pendingRunReposChan is used to communicate repos for which pipeline runs
	// are pending between scheduler and watcher.
	pendingRunReposChan := make(chan string, channelBufferSize)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := &manager.Scheduler{
		TriggeredPipelines: triggeredPipelinesChan,
		TriggeredRepos:     triggeredReposChan,
		PendingRunRepos:    pendingRunReposChan,
		TektonClient:       tClient,
		KubernetesClient:   kClient,
		Logger:             logger.WithTag("scheduler"),
		StorageConfig: manager.StorageConfig{
			Provisioner: storageProvisioner,
			ClassName:   storageClassName,
			Size:        storageSize,
		},
	}
	go s.Run(ctx)

	p := &manager.Pruner{
		TriggeredRepos: triggeredReposChan,
		TektonClient:   tClient,
		Logger:         logger.WithTag("pruner"),
		MinKeepHours:   pruneMinKeepHours,
		MaxKeepRuns:    pruneMaxKeepRuns,
	}
	go p.Run(ctx)

	w := &manager.Watcher{
		PendingRunRepos: pendingRunReposChan,
		Queues:          map[string]bool{},
		TektonClient:    tClient,
		Logger:          logger.WithTag("watcher"),
	}
	go w.Run(ctx)
	// As there is no persistent state, check for queued pipeline runs for all
	// repositories belonging to the Bitbucket project after booting.
	time.AfterFunc(initialWatchWait, func() {
		repos, err := manager.GetRepoNames(bitbucketClient, project)
		if err != nil {
			logger.Warnf("get repo names to check for queued runs: %s", err)
		}
		for _, r := range repos {
			pendingRunReposChan <- r
		}
	})

	r := &manager.BitbucketWebhookReceiver{
		TriggeredPipelines: triggeredPipelinesChan,
		Logger:             logger.WithTag("receiver"),
		BitbucketClient:    bitbucketClient,
		WebhookSecret:      webhookSecret,
		BitbucketWebhookReceiverBase: manager.BitbucketWebhookReceiverBase{
			Namespace: namespace,
			Project:   project,
			RepoBase:  repoBase,
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/health", manager.HealthHandler())
	mux.Handle("/bitbucket", manager.BitbucketHandler(r))
	logger.Infof("Ready to accept requests!")
	return http.ListenAndServe(":8080", mux)
}

func getFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func readIntFromEnvVar(envVar string, fallback int) (int, error) {
	var val int
	valString := os.Getenv(envVar)
	if valString == "" {
		val = fallback
		log.Println(
			"INFO:", envVar, "not set, using default value:", fallback,
		)
	} else {
		i, err := strconv.Atoi(valString)
		if err != nil {
			return 0, fmt.Errorf("could not read value of %s: %s", envVar, err)
		}
		val = i
	}
	return val, nil
}

func readStringFromEnvVar(envVar, fallback string) string {
	val := os.Getenv(envVar)
	if val == "" {
		val = fallback
		log.Printf(
			"INFO: %s not set, using default value: '%s'", envVar, fallback,
		)
	}
	return val
}
