// Package main provides a programm to download all artifacts related to a
// version easily. The artifacts are expected to be in Nexus repositories,
// and placed into the local filesystem. The program not only considers the
// given repository but also any configured subrepositories.
//
// There are two main modes of the program:
// (1) users supply (OpenShift) namespace, (Bitbucket) project, (Git) repository
//     and a tag such as "v1.0.0".
// (2) users run this program from the root of a Git repository and only supply
//     (OpenShift) namespace and tag=WIP. In this case the latest artifacts are
//     downloaded.
// Mode (1) is the main use case, mode (2) is provided as a convenience feature
// for developers.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/installation"
	"github.com/opendevstack/pipeline/internal/repository"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// options configures the program.
// The struct is filled from the flags this program handles.
type options struct {
	kubeconfig      string
	namespace       string
	project         string
	repository      string
	version         bool
	tag             string
	outputDirectory string
	debug           bool
}

// bitbucketArtifactClientInterface is a helper interface that contains the
// methods this program uses on the Bitbucket client. The interface is used
// in testing to mock a Bitbucket client.
type bitbucketArtifactClientInterface interface {
	bitbucket.BranchClientInterface
	bitbucket.TagClientInterface
	bitbucket.RawClientInterface
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
	flag.StringVar(&opts.project, "project", "", "Bitbucket project key of repository")
	flag.StringVar(&opts.repository, "repository", "", "Bitbucket repository key")
	flag.StringVar(&opts.tag, "tag", "", "Git tag to retrieve artifacts for, e.g. v1.0.0 (required)")
	flag.StringVar(&opts.outputDirectory, "output", "artifacts-out", "Directory to place outputs into")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "Enable debug mode")
	flag.BoolVar(&opts.version, "version", false, "Display version of binary")
	flag.Parse()

	if opts.version {
		fmt.Println("0.2.0")
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
	if opts.tag == "" {
		logUsageAndExit("-tag is required")
	}
	if opts.tag != pipelinectxt.WIP {
		if opts.project == "" {
			logUsageAndExit("-project is required when version is not WIP")
		}
		if opts.repository == "" {
			logUsageAndExit("-repository is required when version is not WIP")
		}
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
	nr, err := installation.GetNexusRepositories(c, opts.namespace)
	if err != nil {
		log.Fatalf("Could not get Nexus repositories: %s. Are you logged into the cluster?", err)
	}

	// Bitbucket client
	bcc, err := installation.NewBitbucketClientConfig(c, opts.namespace, logger)
	if err != nil {
		log.Fatalf("Could not create Bitbucket client config: %s. Are you logged into the cluster?", err)
	}
	bitbucketClient := bitbucket.NewClient(bcc)

	err = run(logger, opts, nexusClient, nr, bitbucketClient, workingDir)
	if err != nil {
		log.Fatal(err)
	}
}

// run is the actual main method.
func run(
	logger logging.LeveledLoggerInterface,
	opts options,
	nexusClient nexus.ClientInterface,
	nr *installation.NexusRepositories,
	bitbucketClient bitbucketArtifactClientInterface,
	workingDir string) error {
	// Context
	ctxt, err := getODSContext(opts, bitbucketClient, workingDir)
	if err != nil {
		return err
	}

	err = downloadArtifacts(logger, opts, ctxt, nexusClient, nr)
	if err != nil {
		return err
	}

	// Read ods.yaml file to detect any subrepositories.
	odsConfig, err := getODSConfig(opts, bitbucketClient, workingDir)
	if err != nil {
		return err
	}

	if len(odsConfig.Repositories) > 0 {
		for _, subrepo := range odsConfig.Repositories {
			subrepoCtxt, err := getSubrepoODSContext(ctxt, subrepo, opts, bitbucketClient)
			if err != nil {
				return err
			}
			err = downloadArtifacts(logger, opts, subrepoCtxt, nexusClient, nr)
			if err != nil {
				return err
			}
		}
	}

	return nil
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

// downloadArtifacts downloads the artifact group related to given ODS context.
func downloadArtifacts(
	logger logging.LeveledLoggerInterface,
	opts options,
	ctxt *pipelinectxt.ODSContext,
	nexusClient nexus.ClientInterface,
	nr *installation.NexusRepositories) error {
	artifactsDir := filepath.Join(opts.outputDirectory, opts.tag, ctxt.Repository)
	if _, err := os.Stat(artifactsDir); err == nil {
		return fmt.Errorf("output directory %s already exists", artifactsDir)
	}
	group := pipelinectxt.ArtifactGroupBase(ctxt)
	_, err := pipelinectxt.DownloadGroup(
		nexusClient, []string{nr.Permanent, nr.Temporary}, group, artifactsDir, logger,
	)
	return err
}

// getODSContext assembles an ODS context from given options. If the version is WIP,
// the information is gathered from the Git repository in working directory.
// If the version is not WIP, the information is retrieved from given options
// and the Bitbucket repository.
func getODSContext(opts options, bitbucketClient bitbucket.TagClientInterface, workingDir string) (*pipelinectxt.ODSContext, error) {
	ctxt := &pipelinectxt.ODSContext{
		Namespace: opts.namespace,
	}

	if opts.tag == pipelinectxt.WIP {
		err := ctxt.Assemble(workingDir)
		if err != nil {
			return nil, err
		}
	} else {
		ctxt.Project = opts.project
		ctxt.Repository = opts.repository
		tag, err := bitbucketClient.TagGet(ctxt.Project, ctxt.Repository, opts.tag)
		if err != nil {
			return nil, err
		}
		ctxt.GitCommitSHA = tag.LatestCommit
	}
	return ctxt, nil
}

// getODSConfig reads an ods.y(a)ml file, either from the current directory (if
// tag=WIP) or the remote Bitbucket project identified in the options.
func getODSConfig(opts options, bitbucketClient bitbucket.RawClientInterface, workingDir string) (*config.ODS, error) {
	if opts.tag == pipelinectxt.WIP {
		return config.ReadFromDir(workingDir)
	}
	return repository.GetODSConfig(
		bitbucketClient,
		opts.project,
		opts.repository,
		opts.tag,
	)
}

// getSubrepoODSContext returns an ODS context for the given subrepo.
// The ODS context points to a Git commit, which is either retrieved from the
// best matching branch (if tag=WIP) or the Git tag identified by options.tag.
func getSubrepoODSContext(
	ctxt *pipelinectxt.ODSContext,
	subrepo config.Repository,
	opts options,
	bitbucketClient bitbucketArtifactClientInterface) (*pipelinectxt.ODSContext, error) {
	subrepoCtxt := ctxt.Copy()
	subrepoCtxt.Repository = subrepo.Name
	// For WIP versions, select the best matching branches of subrepositories,
	// and retrieve the latest commit from those branches.
	if opts.tag == pipelinectxt.WIP {
		br, err := repository.BestMatchingBranch(bitbucketClient, subrepoCtxt.Project, subrepo, pipelinectxt.WIP)
		if err != nil {
			return nil, err
		}
		latestCommit, err := repository.LatestCommitForBranch(bitbucketClient, subrepoCtxt.Project, subrepo.Name, br)
		if err != nil {
			return nil, err
		}
		subrepoCtxt.GitCommitSHA = latestCommit
	} else {
		// For non-WIP versions, retrieve the latest commit from given tag.
		tag, err := bitbucketClient.TagGet(subrepoCtxt.Project, subrepoCtxt.Repository, opts.tag)
		if err != nil {
			return nil, err
		}
		subrepoCtxt.GitCommitSHA = tag.LatestCommit
	}
	return subrepoCtxt, nil
}

// logUsageAndExit prints given message followed by flag usage, then exits with exit code 2.
func logUsageAndExit(msg string) {
	log.Println(msg)
	log.Println("Usage:")
	flag.PrintDefaults()
	os.Exit(2)
}
