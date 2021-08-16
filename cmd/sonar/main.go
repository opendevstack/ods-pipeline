package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
)

func main() {
	sonarAuthTokenFlag := flag.String("sonar-auth-token", os.Getenv("SONAR_AUTH_TOKEN"), "sonar-auth-token")
	sonarqubeURLFlag := flag.String("sonar-url", os.Getenv("SONAR_URL"), "sonar-url")
	sonarqubeEditionFlag := flag.String("sonar-edition", os.Getenv("SONAR_EDITION"), "sonar-edition")
	workingDirFlag := flag.String("working-dir", ".", "working directory")
	qualityGateFlag := flag.Bool("quality-gate", false, "require quality gate pass")
	flag.Parse()

	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(".")
	if err != nil {
		log.Fatal(err)
	}
	rootPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(*workingDirFlag)
	if err != nil {
		log.Fatal(err)
	}
	artifactPrefix := ""
	if *workingDirFlag != "." {
		artifactPrefix = strings.Replace(*workingDirFlag, "/", "-", -1)
	}

	sonarClient := sonar.NewClient(&sonar.ClientConfig{
		APIToken:      *sonarAuthTokenFlag,
		BaseURL:       *sonarqubeURLFlag,
		ServerEdition: *sonarqubeEditionFlag,
	})

	// Monorepo support: separate projects in SonarQube.
	// See https://community.sonarsource.com/t/monorepo-and-sonarqube/37990/3.
	sonarProject := fmt.Sprintf("%s-%s", ctxt.Project, ctxt.Component)
	if len(artifactPrefix) > 0 {
		sonarProject = fmt.Sprintf("%s-%s", sonarProject, artifactPrefix)
	}

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
		"author",
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

	if *qualityGateFlag {
		fmt.Println("Checking quality gate ...")
		qualityGateResult, err := sonarClient.QualityGateGet(
			sonar.QualityGateGetParams{Project: sonarProject},
		)
		if err != nil || qualityGateResult.ProjectStatus.Status == "UNKNOWN" {
			fmt.Println("Quality gate status unknown")
			fmt.Println(err)
			os.Exit(1)
		} else if qualityGateResult.ProjectStatus.Status == "ERROR" {
			fmt.Println("Quality gate failed")
			os.Exit(1)
		} else {
			fmt.Println("Quality gate passed")
		}
	}
}
