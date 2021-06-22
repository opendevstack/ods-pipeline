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
	"github.com/opendevstack/pipeline/internal/random"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	kclient "k8s.io/client-go/kubernetes"
)

// SetupFakeRepo writes .ods cache with fake data, without actually initializing a Git repo.
func SetupFakeRepo(t *testing.T, ns, wsDir string) *pipelinectxt.ODSContext {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}
	defer os.Chdir(cwd)
	os.Chdir(wsDir)

	ctxt := &pipelinectxt.ODSContext{
		Namespace:    ns,
		Project:      "myproject",
		Repository:   "myrepo",
		Component:    "myrepo",
		GitCommitSHA: random.PseudoSHA(),
		GitFullRef:   "refs/heads/master",
		GitRef:       "master",
		GitURL:       "http://bitbucket.acme.org/scm/myproject/myrepo.git",
	}
	err = ctxt.WriteCache(wsDir)
	if err != nil {
		t.Fatalf("could not write .ods: %s", err)
	}
	return ctxt
}

// SetupGitRepo initializes a Git repo, commits and writes the result to the .ods cache.
func SetupGitRepo(t *testing.T, ns, wsDir string) *pipelinectxt.ODSContext {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}
	defer os.Chdir(cwd)
	os.Chdir(wsDir)

	initAndCommitOrFatal(t, wsDir)

	ctxt := &pipelinectxt.ODSContext{
		Namespace: ns,
		Project:   "myproject",
		GitURL:    "http://bitbucket.acme.org/scm/myproject/myrepo.git",
	}
	err = ctxt.Assemble(wsDir)
	if err != nil {
		t.Fatalf("could not assemble ODS context information: %s", err)
	}

	err = ctxt.WriteCache(wsDir)
	if err != nil {
		t.Fatalf("could not write .ods: %s", err)
	}
	return ctxt
}

// SetupBitbucketRepo initializes a Git repo, commits, pushes to Bitbucket and writes the result to the .ods cache.
func SetupBitbucketRepo(t *testing.T, c *kclient.Clientset, ns, wsDir, projectKey string) *pipelinectxt.ODSContext {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}
	defer os.Chdir(cwd)
	os.Chdir(wsDir)

	initAndCommitOrFatal(t, wsDir)
	originURL := pushToBitbucketOrFatal(t, c, ns, wsDir, projectKey)

	ctxt := &pipelinectxt.ODSContext{
		Namespace: ns,
		Project:   projectKey,
		GitURL:    originURL,
	}
	err = ctxt.Assemble(wsDir)
	if err != nil {
		t.Fatalf("could not assemble ODS context information: %s", err)
	}

	err = ctxt.WriteCache(wsDir)
	if err != nil {
		t.Fatalf("could not write .ods: %s", err)
	}
	return ctxt
}

func initAndCommitOrFatal(t *testing.T, wsDir string) {
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
	stdout, stderr, err := command.Run("git", []string{"init"})
	if err != nil {
		t.Fatalf("error running git init: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.Run("git", []string{"config", "user.email", "testing@opendevstack.org"})
	if err != nil {
		t.Fatalf("error running git config.user.email: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.Run("git", []string{"config", "user.name", "testing"})
	if err != nil {
		t.Fatalf("error running git config.user.name: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.Run("git", []string{"add", "."})
	if err != nil {
		t.Fatalf("error running git add: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.Run("git", []string{"commit", "-m", "initial commit"})
	if err != nil {
		t.Fatalf("error running git commit: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
}

func pushToBitbucketOrFatal(t *testing.T, c *kclient.Clientset, ns, wsDir, projectKey string) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}
	defer os.Chdir(cwd)
	os.Chdir(wsDir)
	repoName := filepath.Base(wsDir)
	bbURL, err := kubernetes.GetConfigMapKey(c, ns, "ods-bitbucket", "url")
	if err != nil {
		t.Fatalf("could not get Bitbucket URL: %s", err)
	}
	bbURL = "http://localhost:7990"
	bbToken, err := kubernetes.GetSecretKey(c, ns, "ods-bitbucket-auth", "password")
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

func writeFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}
