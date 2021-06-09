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
	subdirectoryFlag := flag.String("subdirectory", ".", "subdirectory to checkout into")
	httpProxyFlag := flag.String("http-proxy", ".", "HTTP_PROXY")
	httpsProxyFlag := flag.String("https-proxy", ".", "HTTPS_PROXY")
	noProxyFlag := flag.String("no-proxy", ".", "NO_PROXY")
	urlFlag := flag.String("url", ".", "URL to clone")
	gitRefFlag := flag.String("git-ref", "", "Git ref to clone")
	sslVerifyFlag := flag.String("ssl-verify", "true", "defines if http.sslVerify should be set to true or false in the global git config")
	submodulesFlag := flag.String("submodules", "true", "defines if the resource should initialize and fetch the submodules")
	depthFlag := flag.String("depth", "1", "performs a shallow clone where only the most recent commit(s) will be fetched")
	flag.Parse()

	// Calculate checkout dir
	// nothing to do, right?
	checkoutDir := *subdirectoryFlag

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

	// git-init
	stdout, stderr, err := command.Run("/ko-app/git-init", []string{
		"-url", *urlFlag,
		"-revision", *gitRefFlag,
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
	// TODO: should we read them before parsing flags and have them as a default?
	// TODO: git ref param: full or short?
	sha, err := getCommitSHA()
	if err != nil {
		log.Fatal(err)
	}
	ctxt := &pipelinectxt.ODSContext{
		Namespace:       *namespaceFlag,
		Project:         *projectFlag,
		Repository:      *repositoryFlag,
		Component:       *componentFlag,
		GitRef:          *gitRefFlag,
		GitFullRef:      *gitRefSpecFlag, // TODO: this is incorrect
		GitCommitSHA:    sha,
		PullRequestBase: *prBaseFlag,
		PullRequestKey:  *prKeyFlag,
	}
	err = ctxt.Assemble(".")
	if err != nil {
		log.Fatal(err)
	}
	err = ctxt.WriteCache(".")
	if err != nil {
		log.Fatal(err)
	}

	// Set Bitbucket build status to "in progress"
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		Timeout:    10 * time.Second,
		APIToken:   *bitbucketAccessTokenFlag,
		MaxRetries: 2,
		BaseURL:    *bitbucketURLFlag,
	})
	url := "http://foo"
	_, _, err = bitbucketClient.BuildStatusPost(ctxt.GitCommitSHA, bitbucket.BuildStatusPostPayload{
		State:       "INPROGRESS",
		Key:         ctxt.GitCommitSHA,
		Name:        ctxt.GitCommitSHA,
		URL:         url,
		Description: "ODS Pipeline Build",
	})
	if err != nil {
		log.Fatal(err)
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
