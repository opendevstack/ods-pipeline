package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/pkg/sonar"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
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

	var namespace string
	if len(*namespaceFlag) > 0 {
		namespace = *namespaceFlag
	} else {
		kubernetesNamespace, err := getTrimmedFileContent(namespaceFile)
		if err != nil {
			panic(err.Error())
		}
		namespace = kubernetesNamespace
	}
	fmt.Printf("namespace=%s\n", namespace)

	var project string
	if len(*projectFlag) > 0 {
		project = *projectFlag
	} else {
		project = strings.TrimSuffix(namespace, "-cd")
	}
	fmt.Printf("project=%s\n", project)

	var repository string
	if len(*repositoryFlag) > 0 {
		repository = *repositoryFlag
	} else {
		r, err := getGitRepository()
		check(err)
		repository = r
	}
	fmt.Printf("repository=%s\n", repository)

	var component string
	if len(*componentFlag) > 0 {
		component = *componentFlag
	} else {
		component = strings.TrimPrefix(repository, fmt.Sprintf("%s-", project))
	}
	fmt.Printf("component=%s\n", component)

	var gitRefSpec string
	if len(*gitRefSpecFlag) > 0 {
		gitRefSpec = *gitRefSpecFlag
	} else {
		grs, err := getGitFullRef()
		check(err)
		gitRefSpec = grs
	}
	fmt.Printf("gitRefSpec=%s\n", gitRefSpec)

	var gitCommitSHA string
	if len(*gitCommitSHAFlag) > 0 {
		gitCommitSHA = *gitCommitSHAFlag
	} else {
		gcs, err := getGitCommitSHA(gitRefSpec)
		check(err)
		gitCommitSHA = gcs
	}
	fmt.Printf("gitCommitSHA=%s\n", gitCommitSHA)
	gitRef, _ := getGitFullRef()

	var prKey string
	if len(*prKeyFlag) > 0 {
		prKey = *prKeyFlag
	} else {
		k, err := getTrimmedFileContent(".ods/pr-key")
		check(err)
		prKey = k
	}
	fmt.Printf("prKey=%s\n", prKey)

	var prBase string
	if len(*prBaseFlag) > 0 {
		prBase = *prBaseFlag
	} else {
		k, err := getTrimmedFileContent(".ods/pr-base")
		check(err)
		prBase = k
	}
	fmt.Printf("prBase=%s\n", prBase)

	sonarClient := sonar.NewClient(&sonar.ClientConfig{
		Timeout:    10 * time.Second,
		APIToken:   *sonarAuthTokenFlag,
		MaxRetries: 2,
		BaseURL:    *sonarqubeURLFlag,
	})

	sonarProject := fmt.Sprintf("%s-%s", project, component)

	fmt.Println("scanning with sonar ...")
	var prInfo *sonar.PullRequest
	if len(prKey) > 0 && prKey != "0" && len(prBase) > 0 {
		prInfo = &sonar.PullRequest{Key: prKey, Branch: gitRef, Base: prBase}
	}
	stdout, err := sonarClient.Scan(
		sonarProject,
		gitRef,
		gitCommitSHA,
		&sonar.BitbucketServer{
			URL:        *bitbucketURLFlag,
			Token:      *bitbucketAccessTokenFlag,
			Project:    project,
			Repository: repository,
		},
		prInfo,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(stdout)

	fmt.Println("generating report ...")
	stdout, err = sonarClient.GenerateReport(sonarProject, "author", gitRef)
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

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func getGitFullRef() (string, error) {
	return getTrimmedFileContent(".ods/git-full-ref")
}

func getGitCommitSHA(refSpec string) (string, error) {
	return getTrimmedFileContent(".ods/git-commit-sha")
}

func getGitCommitSHAInDir(refSpec string, dir string) (string, error) {
	return getTrimmedFileContent(dir + "/.ods/git-commit-sha")
}

func getGitRepository() (string, error) {
	return getTrimmedFileContent(".ods/repository")
}
