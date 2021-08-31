package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
)

type options struct {
	sonarAuthToken string
	sonarURL       string
	sonarEdition   string
	workingDir     string
	qualityGate    bool
	debug          bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.sonarAuthToken, "sonar-auth-token", os.Getenv("SONAR_AUTH_TOKEN"), "sonar-auth-token")
	flag.StringVar(&opts.sonarURL, "sonar-url", os.Getenv("SONAR_URL"), "sonar-url")
	flag.StringVar(&opts.sonarEdition, "sonar-edition", os.Getenv("SONAR_EDITION"), "sonar-edition")
	flag.StringVar(&opts.workingDir, "working-dir", ".", "working directory")
	flag.BoolVar(&opts.qualityGate, "quality-gate", false, "require quality gate pass")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "debug mode")
	flag.Parse()

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	}

	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(".")
	if err != nil {
		log.Fatal(err)
	}
	rootPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(opts.workingDir)
	if err != nil {
		log.Fatal(err)
	}
	artifactPrefix := ""
	if opts.workingDir != "." {
		artifactPrefix = strings.Replace(opts.workingDir, "/", "-", -1) + "-"
	}

	sonarClient := sonar.NewClient(&sonar.ClientConfig{
		APIToken:      opts.sonarAuthToken,
		BaseURL:       opts.sonarURL,
		ServerEdition: opts.sonarEdition,
		Debug:         opts.debug,
		Logger:        logger,
	})

	sonarProject := sonar.ProjectKey(ctxt, artifactPrefix)

	fmt.Println("Scanning with sonar-scanner ...")
	var prInfo *sonar.PullRequest
	if len(ctxt.PullRequestKey) > 0 && ctxt.PullRequestKey != "0" && len(ctxt.PullRequestBase) > 0 {
		prInfo = &sonar.PullRequest{
			Key:    ctxt.PullRequestKey,
			Branch: ctxt.GitRef,
			Base:   ctxt.PullRequestBase,
		}
	}
	stdout, err := sonarClient.Scan(
		sonarProject,
		ctxt.GitRef,
		ctxt.GitCommitSHA,
		prInfo,
	)
	if err != nil {
		fmt.Println(stdout)
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(stdout)

	fmt.Println("Generating reports ...")
	stdout, err = sonarClient.GenerateReports(
		sonarProject,
		"OpenDevStack",
		ctxt.GitRef,
		rootPath,
		artifactPrefix,
	)
	if err != nil {
		fmt.Println(stdout)
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(stdout)

	if opts.qualityGate {
		fmt.Println("Checking quality gate ...")
		qualityGateResult, err := sonarClient.QualityGateGet(
			sonar.QualityGateGetParams{Project: sonarProject},
		)
		if err != nil {
			log.Fatalln(err)
		}
		actualStatus := qualityGateResult.ProjectStatus.Status
		if actualStatus != sonar.QualityGateStatusOk {
			log.Fatalf(
				"Quality gate status is '%s', not '%s'\n",
				actualStatus, sonar.QualityGateStatusOk,
			)
		} else {
			fmt.Println("Quality gate passed")
		}
	}
}
