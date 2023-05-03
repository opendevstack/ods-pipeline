// Package main provides a programm to download all artifacts related to a
// revision easily. The artifacts are expected to be in Nexus repositories,
// and placed into the local filesystem. The program also collects artifacts
// from any subrepositories that were part of the pipeline run producing the
// artifacts.
//
// Run this program from the root of a Git repository and supply the OpenShift
// namespace via -namespace. Example:
//
// ./artifact-download -namespace foo-cd
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/installation"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Version information injected at build time.
var Version string
var GitCommit string

// options configures the program.
// The struct is filled from the flags this program handles.
type options struct {
	kubeconfig      string
	namespace       string
	artifactSource  string
	outputDirectory string
	privateCert     string
	debug           bool
	version         bool
}

func main() {
	workingDir := "."

	kcDefault, err := kubeconfigDefault()
	if err != nil {
		log.Fatal(err)
	}

	opts := options{}
	flag.StringVar(&opts.kubeconfig, "kubeconfig", kcDefault, "Path to kube config file")
	flag.StringVar(&opts.namespace, "namespace", "", "Namespace of ods-pipeline user installation (required)")
	flag.StringVar(&opts.artifactSource, "artifact-source", "", "Artifact source repository")
	flag.StringVar(&opts.outputDirectory, "output", "artifacts-out", "Directory to place outputs into")
	flag.StringVar(&opts.privateCert, "private-cert", "", "Path to private certification (in PEM format)")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "Enable debug mode")
	flag.BoolVar(&opts.version, "version", false, "Display version of binary")
	flag.Parse()

	if opts.version {
		fmt.Println("Version:", Version)
		fmt.Println("Commit: ", GitCommit)
		os.Exit(0)
	}

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	// Validate flags
	if opts.kubeconfig == "" {
		logUsageAndExit("-kubeconfig is required")
	}
	if opts.namespace == "" {
		logUsageAndExit("-namespace is required")
	}
	if opts.artifactSource == "" {
		logUsageAndExit("-artifact-source is required")
	}

	// Kubernetes client
	c, err := newKubernetesClientset(opts.kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Nexus client
	ncc, err := installation.NewNexusClientConfig(c, opts.namespace, logger)
	if err != nil {
		log.Fatalf("Could not create Nexus client config: %s. Are you logged into the cluster?", err)
	}
	nexusClient, err := nexus.NewClient(ncc)
	if err != nil {
		log.Fatal(err)
	}

	err = run(logger, opts, nexusClient, opts.artifactSource, workingDir)
	if err != nil {
		log.Fatal(err)
	}
}

// kubeconfigDefault returns the default location of the Kubernetes config file.
func kubeconfigDefault() (string, error) {
	var kubeconfigDefault string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if homeDir != "" {
		kubeconfigDefault = filepath.Join(homeDir, ".kube", "config")
	}
	return kubeconfigDefault, nil
}

// newKubernetesClientset creates a new Kubernetes clientset from given
// kubeconfig.
func newKubernetesClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("could not build config: %w", err)
	}

	// create the Kubernetes clientset
	kubernetesClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for config: %w", err)
	}

	return kubernetesClientset, nil
}

// logUsageAndExit prints given message followed by flag usage, then exits with exit code 2.
func logUsageAndExit(msg string) {
	log.Println(msg)
	log.Println("Usage:")
	flag.PrintDefaults()
	os.Exit(2)
}
