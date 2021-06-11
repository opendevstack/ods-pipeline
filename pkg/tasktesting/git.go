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

func InitAndCommitOrFatal(t *testing.T, wsDir string) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}
	defer os.Chdir(cwd)
	os.Chdir(wsDir)
	if _, err := os.Stat(".ods"); os.IsNotExist(err) {
		err = os.Mkdir(".ods", 0755)
		if err != nil {
			t.Fatalf("could not create .ods: %s", err)
		}
	}
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
	stdout, stderr, err := command.Run("git", []string{"commit", "-m", "initial commit"})
	if err != nil {
		t.Fatalf("error running git commit: %s, stdout: %s, stderr: %s", err, string(stdout), string(stderr))
	}
}

func PushToBitbucketOrFatal(t *testing.T, c *kclient.Clientset, ns, wsDir, projectKey string) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}
	defer os.Chdir(cwd)
	os.Chdir(wsDir)
	repoName := filepath.Base(wsDir)
	bbURL, err := kubernetes.GetConfigMapKey(c, ns, "bitbucket", "url")
	if err != nil {
		t.Fatalf("could not get Bitbucket URL: %s", err)
	}
	bbURL = "http://localhost:7990"
	bbToken, err := kubernetes.GetSecretKey(c, ns, "bitbucket-auth", "password")
	if err != nil {
		t.Fatalf("could not get Bitbucket token: %s", err)
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
		t.Fatalf("could not create Bitbucket repository: %s", err)
	}

	originURL := fmt.Sprintf("%s/scm/%s/%s.git", bbURL, proj.Key, repo.Slug)

	originURLWithCredentials := strings.Replace(
		originURL,
		"http://",
		fmt.Sprintf("http://%s:%s@", "admin", bbToken),
		-1,
	)
	_, stderr, err := command.Run("git", []string{"remote", "add", "origin", originURLWithCredentials})
	if err != nil {
		t.Fatalf("failed to add remote origin=%s: %s, stderr: %s", originURL, err, stderr)
	}
	_, stderr, err = command.Run("git", []string{"push", "-u", "origin", "master"})
	if err != nil {
		t.Fatalf("failed to push to remote: %s, stderr: %s", err, stderr)
	}

	originURLWithKind := strings.Replace(
		originURL,
		"http://localhost",
		"http://bitbucket-server-test.kind",
		-1,
	)
	return originURLWithKind
}

func WriteDotOdsOrFatal(t *testing.T, wsDir string, projectKey string) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}
	defer os.Chdir(cwd)
	wsName := filepath.Base(wsDir)
	os.Chdir(wsDir)
	err = writeFile(".ods/project", projectKey)
	if err != nil {
		t.Fatalf("could not write .ods/project: %s", err)
	}
	err = writeFile(".ods/repository", wsName)
	if err != nil {
		t.Fatalf("could not write .ods/repository: %s", err)
	}
	err = writeFile(".ods/component", wsName)
	if err != nil {
		t.Fatalf("could not write .ods/component: %s", err)
	}
	sha, err := getTrimmedFileContent(".git/refs/heads/master")
	if err != nil {
		t.Fatalf("error reading .git/refs/heads/master: %s", err)
	}
	err = writeFile(".ods/git-commit-sha", sha)
	if err != nil {
		t.Fatalf("could not write .ods/git-commit-sha: %s", err)
	}
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
