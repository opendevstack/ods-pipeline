package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
)

type options struct {
	sonarAuthToken     string
	sonarURL           string
	sonarEdition       string
	workingDir         string
	rootPath           string
	qualityGate        bool
	trustStore         string
	trustStorePassword string
	debug              bool
}

var defaultOptions = options{
	sonarAuthToken:     os.Getenv("SONAR_AUTH_TOKEN"),
	sonarURL:           os.Getenv("SONAR_URL"),
	sonarEdition:       os.Getenv("SONAR_EDITION"),
	workingDir:         ".",
	qualityGate:        false,
	trustStore:         "${JAVA_HOME}/lib/security/cacerts",
	trustStorePassword: "changeit",
	debug:              (os.Getenv("DEBUG") == "true"),
}

func main() {
	rootPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}

	opts := options{rootPath: rootPath}
	flag.StringVar(&opts.sonarAuthToken, "sonar-auth-token", defaultOptions.sonarAuthToken, "sonar-auth-token")
	flag.StringVar(&opts.sonarURL, "sonar-url", defaultOptions.sonarURL, "sonar-url")
	flag.StringVar(&opts.sonarEdition, "sonar-edition", defaultOptions.sonarEdition, "sonar-edition")
	flag.StringVar(&opts.workingDir, "working-dir", defaultOptions.workingDir, "working directory")
	flag.BoolVar(&opts.qualityGate, "quality-gate", defaultOptions.qualityGate, "require quality gate pass")
	flag.StringVar(&opts.trustStore, "truststore", defaultOptions.trustStore, "JKS truststore")
	flag.StringVar(&opts.trustStorePassword, "truststore-pass", defaultOptions.trustStorePassword, "JKS truststore password")
	flag.BoolVar(&opts.debug, "debug", defaultOptions.debug, "debug mode")
	flag.Parse()

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	ctxt := &pipelinectxt.ODSContext{}
	err = ctxt.ReadCache(".")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(opts.workingDir)
	if err != nil {
		log.Fatal(err)
	}

	sonarClient, err := sonar.NewClient(&sonar.ClientConfig{
		APIToken:           opts.sonarAuthToken,
		BaseURL:            opts.sonarURL,
		ServerEdition:      opts.sonarEdition,
		TrustStore:         opts.trustStore,
		TrustStorePassword: opts.trustStorePassword,
		Debug:              opts.debug,
		Logger:             logger,
	})
	if err != nil {
		log.Fatal("sonar client:", err)
	}

	err = sonarScan(logger, opts, ctxt, sonarClient)
	if err != nil {
		log.Fatal(err)
	}
}

func sonarScan(
	logger logging.LeveledLoggerInterface,
	opts options,
	ctxt *pipelinectxt.ODSContext,
	sonarClient sonar.ClientInterface) error {
	artifactPrefix := ""
	if opts.workingDir != "." {
		artifactPrefix = strings.Replace(opts.workingDir, "/", "-", -1) + "-"
	}

	sonarProject := sonar.ProjectKey(ctxt, artifactPrefix)

	logger.Infof("Scanning with sonar-scanner ...")
	var prInfo *sonar.PullRequest
	if len(ctxt.PullRequestKey) > 0 && ctxt.PullRequestKey != "0" && len(ctxt.PullRequestBase) > 0 {
		logger.Infof("Pull request (ID %s) detected.", ctxt.PullRequestKey)
		prInfo = &sonar.PullRequest{
			Key:    ctxt.PullRequestKey,
			Branch: ctxt.GitRef,
			Base:   ctxt.PullRequestBase,
		}
	}
	err := sonarClient.Scan(
		sonarProject,
		ctxt.GitRef,
		ctxt.GitCommitSHA,
		prInfo,
		os.Stdout,
		os.Stdin,
	)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	logger.Infof("Wait until compute engine task finishes ...")
	err = waitUntilComputeEngineTaskIsSuccessful(logger, sonarClient)
	if err != nil {
		return fmt.Errorf("background task did not finish successfully: %w", err)
	}

	if prInfo == nil {
		logger.Infof("Generating reports ...")
		err := sonarClient.GenerateReports(
			sonarProject,
			"OpenDevStack",
			ctxt.GitRef,
			opts.rootPath,
			artifactPrefix,
		)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}
	} else {
		logger.Infof("No reports are generated for pull request scans.")
	}

	if opts.qualityGate {
		logger.Infof("Checking quality gate ...")
		qualityGateResult, err := sonarClient.QualityGateGet(
			sonar.QualityGateGetParams{Project: sonarProject},
		)
		if err != nil {
			return fmt.Errorf("quality gate could not be retrieved: %w", err)
		}
		err = pipelinectxt.WriteJsonArtifact(
			qualityGateResult,
			filepath.Join(opts.rootPath, pipelinectxt.SonarAnalysisPath),
			fmt.Sprintf("%squality-gate.json", artifactPrefix),
		)
		if err != nil {
			return fmt.Errorf("quality gate status could not be stored as an artifact: %w", err)
		}
		actualStatus := qualityGateResult.ProjectStatus.Status
		if actualStatus != sonar.QualityGateStatusOk {
			return fmt.Errorf(
				"quality gate status is '%s', not '%s'",
				actualStatus, sonar.QualityGateStatusOk,
			)
		} else {
			logger.Infof("Quality gate passed.")
		}
	}

	return nil
}

// waitUntilComputeEngineTaskIsSuccessful reads the scanner report file and
// extracts the task ID. It then waits until the corresponding background task
// in SonarQube succeeds. If the tasks fails or the timeout is reached, an
// error is returned.
func waitUntilComputeEngineTaskIsSuccessful(logger logging.LeveledLoggerInterface, sonarClient sonar.ClientInterface) error {
	reportTaskID, err := sonarClient.ExtractComputeEngineTaskID(sonar.ReportTaskFile)
	if err != nil {
		return fmt.Errorf("cannot read task ID: %w", err)
	}
	params := sonar.ComputeEngineTaskGetParams{ID: reportTaskID}
	attempts := 8 // allows for over 4min task runtime
	sleep := time.Second
	for i := 0; i < attempts; i++ {
		logger.Infof("Waiting %s before checking task status ...", sleep)
		time.Sleep(sleep)
		sleep *= 2
		task, err := sonarClient.ComputeEngineTaskGet(params)
		if err != nil {
			logger.Infof("cannot get status of task: %s", err)
			continue
		}
		switch task.Status {
		case sonar.TaskStatusInProgress:
			logger.Infof("Background task %s has not finished yet", reportTaskID)
		case sonar.TaskStatusPending:
			logger.Infof("Background task %s has not started yet", reportTaskID)
		case sonar.TaskStatusFailed:
			return fmt.Errorf("background task %s has failed", reportTaskID)
		case sonar.TaskStatusSuccess:
			logger.Infof("Background task %s has finished successfully", reportTaskID)
			return nil
		default:
			logger.Infof("Background task %s has unknown status %s", reportTaskID, task.Status)
		}
	}
	return fmt.Errorf("background task %s did not succeed within timeout", reportTaskID)
}
