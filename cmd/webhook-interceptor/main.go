package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/internal/interceptor"
	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
)

const (
	namespaceFile    = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	namespaceSuffix  = "-cd"
	repoBaseEnvVar   = "REPO_BASE"
	tokenEnvVar      = "ACCESS_TOKEN"
	taskKindEnvVar   = "ODS_TASK_KIND"
	taskKindDefault  = "ClusterTask"
	taskSuffixEnvVar = "ODS_TASK_SUFFIX"
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
	if len(repoBase) == 0 {
		return fmt.Errorf("%s must be set", repoBaseEnvVar)
	}

	token := os.Getenv(tokenEnvVar)
	if len(token) == 0 {
		return fmt.Errorf("%s must be set", tokenEnvVar)
	}

	taskKind := os.Getenv(taskKindEnvVar)
	if len(taskKind) == 0 {
		taskKind = taskKindDefault
		log.Println(
			"INFO:",
			taskKindEnvVar,
			"not set, using default value:",
			taskKindDefault,
		)
	}

	taskSuffix := os.Getenv(taskSuffixEnvVar)
	if len(taskSuffix) == 0 {
		log.Println(
			"INFO:",
			taskSuffixEnvVar,
			"not set, using no suffix",
		)
	}

	namespace, err := getFileContent(namespaceFile)
	if err != nil {
		return err
	}

	project := strings.TrimSuffix(namespace, namespaceSuffix)

	// Initialize Tekton client.
	client, err := tektonClient.NewInClusterClient(&tektonClient.ClientConfig{
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

	server, err := interceptor.NewServer(interceptor.ServerConfig{
		Namespace:       namespace,
		Project:         project,
		RepoBase:        repoBase,
		Token:           token,
		TaskKind:        taskKind,
		TaskSuffix:      taskSuffix,
		BitbucketClient: bitbucketClient,
		TektonClient:    client,
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
