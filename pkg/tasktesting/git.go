package tasktesting

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	kclient "k8s.io/client-go/kubernetes"
)

func InitAndCommitOrFatal(t *testing.T, wsDir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}
	defer os.Chdir(cwd)
	os.Chdir(wsDir)
	err = writeFile(".gitignore", ".ods/")
	if err != nil {
		t.Fatalf("could not write .gitignore: %s", err)
	}
	_, stderr, err := command.Run("git", []string{"init"})
	if err != nil {
		t.Fatalf("error running git init: %s, stderr: %s", err, stderr)
	}
	_, stderr, err = command.Run("git", []string{"config", "user.email", "testing@opendevstack.org"})
	if err != nil {
		t.Fatalf("error running git config.user.email: %s, stderr: %s", err, stderr)
	}
	_, stderr, err = command.Run("git", []string{"config", "user.name", "testing"})
	if err != nil {
		t.Fatalf("error running git config.user.name: %s, stderr: %s", err, stderr)
	}
	_, stderr, err = command.Run("git", []string{"add", "."})
	if err != nil {
		t.Fatalf("error running git add: %s, stderr: %s", err, stderr)
	}
	_, stderr, err = command.Run("git", []string{"commit", "-m", "initial commit"})
	if err != nil {
		t.Fatalf("error running git commit: %s, stderr: %s", err, stderr)
	}
	return nil
}

func PushToBitbucket(c *kclient.Clientset, ns string, projectKey string, repoName string) error {
	bbURL, err := kubernetes.GetConfigMapKey(c, ns, "bitbucket", "url")
	if err != nil {
		return err
	}
	bbURL = "http://localhost:7990"
	bbToken, err := kubernetes.GetSecretKey(c, ns, "bitbucket-auth", "password")
	if err != nil {
		return err
	}

	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		Timeout:    10 * time.Second,
		APIToken:   bbToken,
		MaxRetries: 2,
		BaseURL:    bbURL,
		Logger:     &logging.LeveledLogger{Level: logging.LevelDebug},
	})

	proj := bitbucket.Project{Key: projectKey}
	repo, err := bitbucketClient.RepoCreate(proj.Key, bitbucket.RepoCreatePayload{
		Name:          repoName,
		SCMID:         "git",
		Forkable:      true,
		DefaultBranch: "master",
	})
	if err != nil {
		return err
	}

	bbCredentialsURL := strings.Replace(
		bbURL,
		"http://",
		fmt.Sprintf("http://%s:%s@", "admin", bbToken),
		-1,
	)
	origin := fmt.Sprintf("%s/scm/%s/%s.git", bbCredentialsURL, proj.Key, repo.Slug)
	_, stderr, err := command.Run("git", []string{"remote", "add", "origin", origin})
	if err != nil {
		return fmt.Errorf("failed to add remote origin=%s: %s, stderr: %s", origin, err, stderr)
	}
	_, stderr, err = command.Run("git", []string{"push", "-u", "origin", "master"})
	if err != nil {
		return fmt.Errorf("failed to push to remote: %s, stderr: %s", err, stderr)
	}
	return nil
}

func WriteDotOds(wsDir string, projectKey string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)
	wsName := filepath.Base(wsDir)
	os.Chdir(wsDir)
	err = writeFile(".ods/project", projectKey)
	if err != nil {
		return err
	}
	err = writeFile(".ods/repository", wsName)
	if err != nil {
		return err
	}
	err = writeFile(".ods/component", wsName)
	if err != nil {
		return err
	}
	sha, err := getTrimmedFileContent(".git/refs/heads/master")
	if err != nil {
		return err
	}
	err = writeFile(".ods/git-commit-sha", sha)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
