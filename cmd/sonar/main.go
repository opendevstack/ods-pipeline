package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
)

func main() {
	sonarAuthTokenFlag := flag.String("sonar-auth-token", os.Getenv("SONAR_AUTH_TOKEN"), "sonar-auth-token")
	sonarqubeURLFlag := flag.String("sonar-url", os.Getenv("SONAR_URL"), "sonar-url")
	sonarqubeEditionFlag := flag.String("sonar-edition", os.Getenv("SONAR_EDITION"), "sonar-edition")
	qualityGateFlag := flag.Bool("quality-gate", false, "require quality gate pass")
	flag.Parse()

	ctxt := &pipelinectxt.ODSContext{}
	err := ctxt.ReadCache(".")
	if err != nil {
		panic(err.Error())
	}

	sonarClient := sonar.NewClient(&sonar.ClientConfig{
		APIToken:      *sonarAuthTokenFlag,
		BaseURL:       *sonarqubeURLFlag,
		ServerEdition: *sonarqubeEditionFlag,
	})

	sonarProject := fmt.Sprintf("%s-%s", ctxt.Project, ctxt.Component)

	fmt.Println("scanning with sonar ...")
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
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(stdout)

	fmt.Println("generating report ...")
	stdout, err = sonarClient.GenerateReport(sonarProject, "author", ctxt.GitRef)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(stdout)

	if *qualityGateFlag {
		fmt.Println("checking quality gate ...")
		qualityGateResult, err := sonarClient.QualityGateGet(
			sonar.QualityGateGetParams{Project: sonarProject},
		)
		if err != nil || qualityGateResult.ProjectStatus.Status == "UNKNOWN" {
			fmt.Println("quality gate unknown")
			fmt.Println(err)
			os.Exit(1)
		} else if qualityGateResult.ProjectStatus.Status == "ERROR" {
			fmt.Println("quality gate failed")
			os.Exit(1)
		} else {
			fmt.Println("quality gate passed")
		}
	}
}
