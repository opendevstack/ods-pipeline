package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/logging"
)

// gitCheckoutParams holds the parameters configuring the checkout.
type gitCheckoutParams struct {
	repoURL              string
	bitbucketAccessToken string
	recurseSubmodules    string
	depth                string
	fullRef              string
}

// gitCheckout encapsulates the steps required to perform a Git checkout
func gitCheckout(p gitCheckoutParams) (err error) {
	fetchArgs := []string{"--recurse-submodules=" + p.recurseSubmodules}
	if p.depth != "" {
		fetchArgs = append(fetchArgs, "--depth="+p.depth)
	}
	fetchArgs = append(fetchArgs, "origin", "--update-head-ok", "--force", p.fullRef)
	steps := [][]string{
		{"init"},
		// Even though Tekton prepares credentials to be used for each task,
		// we set the auth explicitly here. The motivation is that Tekton uses
		// basic auth to pass the username/token, which fails in environments
		// that have basic auth disabled for Bitbucket.
		{"config",
			fmt.Sprintf("http.%s.extraHeader", p.repoURL),
			fmt.Sprintf("Authorization: Bearer %s", p.bitbucketAccessToken),
		},
		{"config",
			fmt.Sprintf("http.%s/info/lfs.extraHeader", p.repoURL),
			fmt.Sprintf("Authorization: Bearer %s", p.bitbucketAccessToken),
		},
		{"remote", "add", "origin", p.repoURL},
		append([]string{"fetch"}, fetchArgs...),
		{"checkout", "-f", "FETCH_HEAD"},
	}
	for _, args := range steps {
		if err == nil {
			err = runGitCmd(args...)
		}
	}
	return
}

// runGitCmd executes git with given args.
func runGitCmd(args ...string) error {
	var output bytes.Buffer
	err := command.Run("git", args, []string{}, &output, &output)
	if err != nil {
		return fmt.Errorf("git %v: %w\n%s", args, err, output.String())
	}
	return nil
}

func getCommitSHA(dir string) (string, error) {
	content, err := os.ReadFile(filepath.Join(dir, ".git/HEAD"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func gitLfsInUse(logger logging.LeveledLoggerInterface, dir string) (lfs bool, err error) {
	stdout, stderr, err := command.RunBufferedInDir("git", []string{"lfs", "ls-files", "--all"}, dir)
	if err != nil {
		return false, fmt.Errorf("cannot list git lfs files: %s (%w)", stderr, err)
	}
	return strings.TrimSpace(string(stdout)) != "", err
}

func gitLfsEnableAndPullFiles(logger logging.LeveledLoggerInterface, dir string) (err error) {
	stdout, stderr, err := command.RunBufferedInDir("git", []string{"lfs", "install"}, dir)
	if err != nil {
		return fmt.Errorf("lfs install: %s (%w)", stderr, err)
	}
	logger.Infof(string(stdout))
	stdout, stderr, err = command.RunBufferedInDir("git", []string{"lfs", "pull"}, dir)
	if err != nil {
		return fmt.Errorf("lfs pull: %s (%w)", stderr, err)
	}
	logger.Infof(string(stdout))
	return err
}
