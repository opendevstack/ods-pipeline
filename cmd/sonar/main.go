package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
)

func main() {
	sonarAuthTokenFlag := flag.String("sonar-auth-token", os.Getenv("SONAR_AUTH_TOKEN"), "sonar-auth-token")
	sonarqubeURLFlag := flag.String("sonar-url", os.Getenv("SONAR_URL"), "sonar-url")
	qualityGateFlag := flag.Bool("quality-gate", false, "require quality gate pass")

	bitbucketAccessTokenFlag := flag.String("bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	bitbucketURLFlag := flag.String("bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	namespaceFlag := flag.String("namespace", "", "namespace")
	projectFlag := flag.String("project", "", "project")
	repositoryFlag := flag.String("repository", "", "repository")
	componentFlag := flag.String("component", "", "component")
	gitRefSpecFlag := flag.String("git-ref-spec", "", "Git ref spec")
	gitCommitSHAFlag := flag.String("git-commit-sha", "", "Git commit SHA")
	prKeyFlag := flag.String("pr-key", "", "pull request key")
	prBaseFlag := flag.String("pr-base", "", "pull request base")
	flag.Parse()

	// TODO: should we read them before parsing flags and have them as a default?
	// TODO: git ref param: full or short?
	ctxt := &pipelinectxt.ODSContext{
		Namespace:       *namespaceFlag,
		Project:         *projectFlag,
		Repository:      *repositoryFlag,
		Component:       *componentFlag,
		GitFullRef:      *gitRefSpecFlag,
		GitCommitSHA:    *gitCommitSHAFlag,
		PullRequestBase: *prBaseFlag,
		PullRequestKey:  *prKeyFlag,
	}
	err := ctxt.ReadCache(".")
	if err != nil {
		panic(err.Error())
	}

	sonarClient := sonar.NewClient(&sonar.ClientConfig{
		Timeout:       10 * time.Second,
		APIToken:      *sonarAuthTokenFlag,
		MaxRetries:    2,
		BaseURL:       *sonarqubeURLFlag,
		ServerEdition: "community",
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
		&sonar.BitbucketServer{
			URL:        *bitbucketURLFlag,
			Token:      *bitbucketAccessTokenFlag,
			Project:    ctxt.Project,
			Repository: ctxt.Repository,
		},
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
