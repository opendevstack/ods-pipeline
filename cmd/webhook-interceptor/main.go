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
)

const (
	namespaceFile     = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	tokenFile         = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	caCert            = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	pipelineFilename  = "ods.yml"
	tektonAPIBasePath = "/apis/tekton.dev/v1beta1"
	letterBytes       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	namespaceSuffix   = "-cd"
	apiHostEnvVar     = "API_HOST"
	apiHostDefault    = "openshift.default.svc.cluster.local"
	repoBaseEnvVar    = "REPO_BASE"
	tokenEnvVar       = "ACCESS_TOKEN"
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

	apiHost := os.Getenv(apiHostEnvVar)
	if len(apiHost) == 0 {
		apiHost = apiHostDefault
		log.Println(
			"INFO:",
			apiHostEnvVar,
			"not set, using default value:",
			apiHostDefault,
		)
	}

	namespace, err := getFileContent(namespaceFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	project := strings.TrimSuffix(namespace, namespaceSuffix)

	client, err := interceptor.NewClient(apiHost, namespace)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	server := interceptor.NewServer(client, namespace, project, repoBase, token)

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

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
