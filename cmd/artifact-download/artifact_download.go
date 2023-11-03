package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opendevstack/ods-pipeline/pkg/artifact"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
	"github.com/opendevstack/ods-pipeline/pkg/nexus"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
	"sigs.k8s.io/yaml"
)

// run is the actual main method.
func run(
	logger logging.LeveledLoggerInterface,
	opts options,
	artifactClient nexus.ClientInterface,
	artifactSourceRepository string,
	workingDir string) error {
	ctxt, err := assembleODSContext(opts.namespace, workingDir)
	if err != nil {
		return fmt.Errorf("assemble ODS context: %w", err)
	}

	logger.Infof("Downloading artifacts of repository ...")
	err = downloadArtifacts(
		logger,
		filepath.Join(opts.outputDirectory, ctxt.Repository),
		ctxt,
		artifactClient, artifactSourceRepository,
	)
	if err != nil {
		return fmt.Errorf("download artifacts: %w", err)
	}

	// Handle subrepositories
	artifactsDir := filepath.Join(opts.outputDirectory, ctxt.Repository)
	pra, err := readPipelineRunArtifact(artifactsDir)
	if err != nil {
		return fmt.Errorf("read pipeline run artifact: %w", err)
	}
	if len(pra.Repositories) > 0 {
		logger.Infof("Downloading artifacts of subrepositories ...")
	}
	for subrepoName, subrepoCommit := range pra.Repositories {
		subrepoCtxt := ctxt.Copy()
		subrepoCtxt.Repository = subrepoName
		subrepoCtxt.GitCommitSHA = subrepoCommit
		err = downloadArtifacts(
			logger,
			filepath.Join(opts.outputDirectory, subrepoName),
			subrepoCtxt,
			artifactClient, artifactSourceRepository,
		)
		if err != nil {
			return fmt.Errorf("download artifacts of %s: %w", subrepoName, err)
		}
	}

	return nil
}

func readPipelineRunArtifact(artifactsDir string) (*artifact.PipelineRun, error) {
	searchDir := filepath.Join(artifactsDir, pipelinectxt.PipelineRunsDir)
	filename, err := findPipelineRunArtifact(searchDir)
	if err != nil {
		return nil, fmt.Errorf("find pipeline run artifact: %w", err)
	}

	f, err := os.ReadFile(filepath.Join(searchDir, filename))
	if err != nil {
		return nil, fmt.Errorf("read pipeline run artifact file: %w", err)
	}

	var pra artifact.PipelineRun
	err = yaml.Unmarshal(f, &pra)
	if err != nil {
		return nil, fmt.Errorf("unmarshal pipeline run artifact: %w", err)
	}

	return &pra, err
}

func findPipelineRunArtifact(searchDir string) (string, error) {
	files, err := os.ReadDir(searchDir)
	if err != nil {
		return "", fmt.Errorf("read directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			return file.Name(), nil
		}
	}
	return "", errors.New("no file found in " + searchDir)
}

// downloadArtifacts downloads the artifact group related to given ODS context.
func downloadArtifacts(
	logger logging.LeveledLoggerInterface,
	outputDirectory string,
	ctxt *pipelinectxt.ODSContext,
	client nexus.ClientInterface,
	sourceRepository string) error {
	if _, err := os.Stat(outputDirectory); err == nil {
		return fmt.Errorf("output directory %s already exists", outputDirectory)
	}
	group := pipelinectxt.ArtifactGroupBase(ctxt)
	_, err := pipelinectxt.DownloadGroup(
		client, sourceRepository, group, outputDirectory, logger,
	)
	return err
}

// assembleODSContext assembles an ODS context from given options.
// The information is gathered from the Git repository in working directory.
func assembleODSContext(namespace string, workingDir string) (*pipelinectxt.ODSContext, error) {
	ctxt := &pipelinectxt.ODSContext{
		Namespace: namespace,
	}
	err := ctxt.Assemble(workingDir, "")
	return ctxt, err
}
