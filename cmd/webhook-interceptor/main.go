package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/internal/interceptor"
	kubernetesClient "github.com/opendevstack/pipeline/internal/kubernetes"
	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
)

const (
	namespaceFile            = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	namespaceSuffix          = "-cd"
	repoBaseEnvVar           = "REPO_BASE"
	tokenEnvVar              = "ACCESS_TOKEN"
	taskKindEnvVar           = "ODS_TASK_KIND"
	taskKindDefault          = "ClusterTask"
	taskSuffixEnvVar         = "ODS_TASK_SUFFIX"
	storageProvisionerEnvVar = "ODS_STORAGE_PROVISIONER"
	storageClassNameEnvVar   = "ODS_STORAGE_CLASS_NAME"
	storageClassNameDefault  = "standard"
	storageSizeEnvVar        = "ODS_STORAGE_SIZE"
	storageSizeDefault       = "2Gi"
	pruneMinKeepHoursEnvVar  = "ODS_PRUNE_MIN_KEEP_HOURS"
	pruneMinKeepHoursDefault = 48
	pruneMaxKeepRunsEnvVar   = "ODS_PRUNE_MAX_KEEP_RUNS"
	pruneMaxKeepRunsDefault  = 20
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
	log.Println("Booting")

	repoBase := os.Getenv(repoBaseEnvVar)
	if repoBase == "" {
		return fmt.Errorf("%s must be set", repoBaseEnvVar)
	}

	token := os.Getenv(tokenEnvVar)
	if token == "" {
		return fmt.Errorf("%s must be set", tokenEnvVar)
	}

	taskKind := readStringFromEnvVar(taskKindEnvVar, taskKindDefault)

	taskSuffix := readStringFromEnvVar(taskSuffixEnvVar, "")

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
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: token,
		BaseURL:  strings.TrimSuffix(repoBase, "/scm"),
	})

	// TODO: Use this logger in the interceptor as well, not just in the pruner.
	var logger logging.LeveledLoggerInterface
	if os.Getenv("DEBUG") == "true" {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	pruner, err := interceptor.NewPipelineRunPrunerByStage(
		tClient,
		logger,
		pruneMinKeepHours,
		pruneMaxKeepRuns,
	)
	if err != nil {
		return fmt.Errorf("could not create pruner: %w", err)
	}

	server, err := interceptor.NewServer(interceptor.ServerConfig{
		Namespace:  namespace,
		Project:    project,
		RepoBase:   repoBase,
		Token:      token,
		TaskKind:   taskKind,
		TaskSuffix: taskSuffix,
		StorageConfig: interceptor.StorageConfig{
			Provisioner: storageProvisioner,
			ClassName:   storageClassName,
			Size:        storageSize,
		},
		KubernetesClient:  kClient,
		TektonClient:      tClient,
		BitbucketClient:   bitbucketClient,
		PipelineRunPruner: pruner,
	})
	if err != nil {
		return err
	}

	log.Println("Ready to accept requests")

	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(health))
	mux.Handle("/", http.HandlerFunc(server.HandleRoot))
	return http.ListenAndServe(":8080", mux)
}

func health(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`{"health":"ok"}`))
	if err != nil {
		http.Error(w, `{"health":"error"}`, http.StatusInternalServerError)
		return
	}
}

func getFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
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
