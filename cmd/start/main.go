package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

func main() {
	bitbucketAccessTokenFlag := flag.String("bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	bitbucketURLFlag := flag.String("bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	namespaceFlag := flag.String("namespace", "", "namespace")
	projectFlag := flag.String("project", "", "project")
	repositoryFlag := flag.String("repository", "", "repository")
	componentFlag := flag.String("component", "", "component")
	gitRefSpecFlag := flag.String("git-ref-spec", "", "(optional) git refspec to fetch before checking out revision")
	//gitCommitSHAFlag := flag.String("git-commit-sha", "", "Git commit SHA")
	prKeyFlag := flag.String("pr-key", "", "pull request key")
	prBaseFlag := flag.String("pr-base", "", "pull request base")
	httpProxyFlag := flag.String("http-proxy", ".", "HTTP_PROXY")
	httpsProxyFlag := flag.String("https-proxy", ".", "HTTPS_PROXY")
	noProxyFlag := flag.String("no-proxy", ".", "NO_PROXY")
	urlFlag := flag.String("url", ".", "URL to clone")
	gitFullRefFlag := flag.String("git-full-ref", "", "Git (full) ref to clone")
	sslVerifyFlag := flag.String("ssl-verify", "true", "defines if http.sslVerify should be set to true or false in the global git config")
	submodulesFlag := flag.String("submodules", "true", "defines if the resource should initialize and fetch the submodules")
	depthFlag := flag.String("depth", "1", "performs a shallow clone where only the most recent commit(s) will be fetched")
	consoleURLFlag := flag.String("console-url", "", "web console URL")
	pipelineRunNameFlag := flag.String("pipeline-run-name", "", "name of pipeline run")
	flag.Parse()

	checkoutDir := "."

	// clean dir
	err := deleteDirectoryContents(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}

	// set proxy env vars
	if len(*httpProxyFlag) > 0 {
		err = os.Setenv("HTTP_PROXY", *httpProxyFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(*httpsProxyFlag) > 0 {
		err = os.Setenv("HTTPS_PROXY", *httpsProxyFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(*noProxyFlag) > 0 {
		err = os.Setenv("NO_PROXY", *noProxyFlag)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctxt := prepareODSContextForRepo(checkoutDir, urlFlag, gitFullRefFlag, gitRefSpecFlag, sslVerifyFlag, submodulesFlag, depthFlag, namespaceFlag, projectFlag, repositoryFlag, componentFlag, prBaseFlag, prKeyFlag)

	// Set Bitbucket build status to "in progress"
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		Timeout:    10 * time.Second,
		APIToken:   *bitbucketAccessTokenFlag,
		MaxRetries: 2,
		BaseURL:    *bitbucketURLFlag,
	})
	pipelineRunURL := fmt.Sprintf(
		"%s/k8s/ns/%s/tekton.dev~v1beta1~PipelineRun/%s/",
		*consoleURLFlag,
		ctxt.Namespace,
		*pipelineRunNameFlag,
	)
	err = bitbucketClient.BuildStatusCreate(ctxt.GitCommitSHA, bitbucket.BuildStatusCreatePayload{
		State:       "INPROGRESS",
		Key:         ctxt.GitCommitSHA,
		Name:        ctxt.GitCommitSHA,
		URL:         pipelineRunURL,
		Description: "ODS Pipeline Build",
	})
	if err != nil {
		log.Fatal(err)
	}

	// If ods.yml is present, clone only first-level of child repositories
	odsConfig, err := config.GetODSConfig("ods.yml")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d child repositories found\n", len(odsConfig.Repositories))

	for _, repo := range odsConfig.Repositories {
		log.Printf("Repository name: %s, url: %s\n", repo.Name, repo.URL)
		checkoutDir = fmt.Sprintf(".ods/repos/%s", repo)
		prepareODSContextForRepo(checkoutDir, urlFlag, gitFullRefFlag, gitRefSpecFlag, sslVerifyFlag, submodulesFlag, depthFlag, namespaceFlag, projectFlag, repositoryFlag, componentFlag, prBaseFlag, prKeyFlag)

	}

}

func deleteDirectoryContents(directory string) error {
	// Open the directory and read all its files.
	dirRead, err := os.Open(directory)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", directory, err)
	}
	dirFiles, err := dirRead.Readdir(0)
	if err != nil {
		return fmt.Errorf("could not read files in %s: %w", directory, err)
	}

	// Loop over the directory's files and remove them.
	for _, f := range dirFiles {
		filename := filepath.Join(directory, f.Name())
		err := os.RemoveAll(filename)
		if err != nil {
			return fmt.Errorf("could not remove file %s: %w", filename, err)
		}
	}
	return nil
}

func getCommitSHA() (string, error) {
	content, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func prepareODSContextForRepo(checkoutDir string, urlFlag, gitFullRefFlag, gitRefSpecFlag, sslVerifyFlag, submodulesFlag, depthFlag, namespaceFlag, projectFlag, repositoryFlag, componentFlag, prBaseFlag, prKeyFlag *string) *pipelinectxt.ODSContext {
	// git-init
	stdout, stderr, err := command.Run("/ko-app/git-init", []string{
		"-url", *urlFlag,
		"-revision", *gitFullRefFlag,
		"-refspec", *gitRefSpecFlag,
		"-path", checkoutDir,
		"-sslVerify", *sslVerifyFlag,
		"-submodules", *submodulesFlag,
		"-depth", *depthFlag,
	})
	if err != nil {
		log.Println(string(stderr))
		log.Fatal(err)
	}
	fmt.Println(string(stdout))

	// write ODS cache
	sha, err := getCommitSHA()
	if err != nil {
		log.Fatal(err)
	}
	ctxt := &pipelinectxt.ODSContext{
		Namespace:       *namespaceFlag,
		Project:         *projectFlag,
		Repository:      *repositoryFlag,
		Component:       *componentFlag,
		GitFullRef:      *gitFullRefFlag,
		GitCommitSHA:    sha,
		PullRequestBase: *prBaseFlag,
		PullRequestKey:  *prKeyFlag,
	}
	err = ctxt.Assemble(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}
	err = ctxt.WriteCache(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}

	return ctxt
}
