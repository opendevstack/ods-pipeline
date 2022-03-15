package tasktesting

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/random"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
	kclient "k8s.io/client-go/kubernetes"
)

// SetupFakeRepo writes .ods cache with fake data, without actually initializing a Git repo.
func SetupFakeRepo(t *testing.T, ns, wsDir string) *pipelinectxt.ODSContext {
	ctxt := &pipelinectxt.ODSContext{
		Namespace:    ns,
		Project:      "myproject",
		Repository:   "myrepo",
		Component:    "myrepo",
		GitCommitSHA: random.PseudoSHA(),
		GitFullRef:   "refs/heads/master",
		GitRef:       "master",
		GitURL:       "http://bitbucket.acme.org/scm/myproject/myrepo.git",
		Environment:  "dev",
		Version:      pipelinectxt.WIP,
	}
	err := ctxt.WriteCache(wsDir)
	if err != nil {
		t.Fatalf("could not write %s: %s", pipelinectxt.BaseDir, err)
	}
	return ctxt
}

func assembleAndCacheOdsCtxtOrFatal(t *testing.T, ctxt *pipelinectxt.ODSContext, wsDir string) {
	err := ctxt.Assemble(wsDir)
	if err != nil {
		t.Fatalf("could not assemble ODS context information: %s", err)
	}
	err = ctxt.WriteCache(wsDir)
	if err != nil {
		t.Fatalf("could not write %s: %s", pipelinectxt.BaseDir, err)
	}
}

// SetupGitRepo initializes a Git repo, commits and writes the result to the .ods cache.
func SetupGitRepo(t *testing.T, ns, wsDir string) *pipelinectxt.ODSContext {
	initAndCommitOrFatal(t, wsDir)
	bbURL := "http://localhost:7990"
	repoName := filepath.Base(wsDir)
	ctxt := &pipelinectxt.ODSContext{
		Namespace:   ns,
		GitURL:      fmt.Sprintf("%s/scm/%s/%s.git", bbURL, BitbucketProjectKey, repoName),
		Environment: "dev",
		Version:     pipelinectxt.WIP,
	}
	assembleAndCacheOdsCtxtOrFatal(t, ctxt, wsDir)
	return ctxt
}

// SetupBitbucketRepo initializes a Git repo, commits, pushes to Bitbucket and writes the result to the .ods cache.
func SetupBitbucketRepo(t *testing.T, c *kclient.Clientset, ns, wsDir, projectKey string) *pipelinectxt.ODSContext {
	initAndCommitOrFatal(t, wsDir)
	originURL := pushToBitbucketOrFatal(t, c, ns, wsDir, projectKey)

	ctxt := &pipelinectxt.ODSContext{
		Namespace:   ns,
		Project:     projectKey,
		GitURL:      originURL,
		Environment: "dev",
		Version:     pipelinectxt.WIP,
	}
	assembleAndCacheOdsCtxtOrFatal(t, ctxt, wsDir)
	return ctxt
}

func initAndCommitOrFatal(t *testing.T, wsDir string) {
	if _, err := os.Stat(pipelinectxt.BaseDir); os.IsNotExist(err) {
		err = os.Mkdir(pipelinectxt.BaseDir, 0755)
		if err != nil {
			t.Fatalf("could not create %s: %s", pipelinectxt.BaseDir, err)
		}
	}
	if err := pipelinectxt.WriteGitIgnore(filepath.Join(wsDir, ".gitignore")); err != nil {
		t.Fatalf("could not write .gitignore: %s", err)
	}
	stdout, stderr, err := command.RunInDir("git", []string{"init"}, wsDir)
	if err != nil {
		t.Fatalf("error running git init: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"config", "user.email", "testing@opendevstack.org"}, wsDir)
	if err != nil {
		t.Fatalf("error running git config.user.email: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"config", "user.name", "testing"}, wsDir)
	if err != nil {
		t.Fatalf("error running git config.user.name: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"add", "."}, wsDir)
	if err != nil {
		t.Fatalf("error running git add: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"commit", "-m", "initial commit"}, wsDir)
	if err != nil {
		t.Fatalf("error running git commit: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
}

func PushFileToBitbucketOrFatal(t *testing.T, c *kclient.Clientset, ns, wsDir, branch, filename string) {
	stdout, stderr, err := command.RunInDir("git", []string{"add", filename}, wsDir)
	if err != nil {
		t.Fatalf("failed to add file=%s: %s, stdout: %s, stderr: %s", filename, err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"commit", "-m", "update " + filename}, wsDir)
	if err != nil {
		t.Fatalf("error running git commit: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"push", "origin", branch}, wsDir)
	if err != nil {
		t.Fatalf("failed to push to remote: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
}

func pushToBitbucketOrFatal(t *testing.T, c *kclient.Clientset, ns, wsDir, projectKey string) string {
	repoName := filepath.Base(wsDir)
	bbURL := "http://localhost:7990"
	bbToken, err := kubernetes.GetSecretKey(c, ns, "ods-bitbucket-auth", "password")
	if err != nil {
		t.Fatalf("could not get Bitbucket token: %s", err)
	}

	bitbucketClient := BitbucketClientOrFatal(t, c, ns)

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
	_, stderr, err := command.RunInDir("git", []string{"remote", "add", "origin", originURLWithCredentials}, wsDir)
	if err != nil {
		t.Fatalf("failed to add remote origin=%s: %s, stderr: %s", originURL, err, stderr)
	}
	_, stderr, err = command.RunInDir("git", []string{"push", "-u", "origin", "master"}, wsDir)
	if err != nil {
		t.Fatalf("failed to push to remote: %s, stderr: %s", err, stderr)
	}

	originURLWithKind := strings.Replace(
		originURL,
		"http://localhost",
		"http://ods-test-bitbucket-server.kind",
		-1,
	)
	return originURLWithKind
}

// EnableLfsOnBitbucketRepoOrFatal enable Git LFS extension in existing repo, using private Bitbucket API.
func EnableLfsOnBitbucketRepoOrFatal(t *testing.T, repo, projectKey string) {
	// we need to handle cookies in the http client for the next requests
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("failed to initialize cookie jar, err: %s", err)
	}
	client := &http.Client{
		Jar: jar,
		// we don't want to follow redirects since 200 pages are not useful (can be error page)
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// login, start bitbucket session and get cookies
	requestBody := strings.NewReader("j_username=admin&j_password=admin&_atl_remember_me=on&submit=Login")
	req, err := http.NewRequest(http.MethodPost, "http://localhost:7990/j_atl_security_check", requestBody)
	if err != nil {
		t.Fatalf("failed to prepare login request with cookie jar, err: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to run login request with cookie jar, err: %s", err)
	}
	if res.StatusCode != http.StatusFound {
		t.Fatalf("failed to login with cookie jar, status code was: %d", res.StatusCode)
	}
	// we need to GET the settings page to parse the atl_token required afterwards
	settingsUrl := fmt.Sprintf("http://localhost:7990/projects/%s/repos/%s/settings", projectKey, repo)
	res, err = client.Get(settingsUrl)
	if err != nil {
		t.Fatalf("failed to request settings page with cookie jar, err: %s", err)
	}
	// we are parsing the html site to find the hidden atl_token input form attribute value
	atlToken, err := getAtlassianToken(res.Body)
	if err != nil {
		t.Fatalf("failed to get atl_token attribute, err: %s", err)
	}
	// now we can do a form request to setup the setup the repo enabling LFS
	requestBody = strings.NewReader(fmt.Sprintf("name=%s&description=&defaultBranchId=refs/heads/master&forkable=on&lfsRepoEnabled=on&submit=Save&atl_token=%s", repo, atlToken))
	req, err = http.NewRequest(http.MethodPost, settingsUrl, requestBody)
	if err != nil {
		t.Fatalf("failed to prepare enable LFS request with cookie jar, err: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to run enable LFS request with cookie jar, err: %s", err)
	}
	if res.StatusCode != http.StatusFound {
		t.Fatalf("failed to enable LFS with cookie jar, status code was: %d", res.StatusCode)
	}
}

func getAtlassianToken(b io.ReadCloser) (string, error) {
	var finder func(*html.Node) (string, bool)
	finder = func(n *html.Node) (string, bool) {
		if n.Type == html.ElementNode && n.Data == "input" {
			ourAttr, ok := getAttribute(n, "name")
			if ok && ourAttr == "atl_token" {
				return getAttribute(n, "value")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if val, ok := finder(c); ok {
				return val, ok
			}
		}
		return "", false
	}
	doc, err := html.Parse(b)
	if err != nil {
		return "", err
	}
	if val, ok := finder(doc); ok {
		return val, nil
	}
	return "", errors.New("atl_token not found")
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

// UpdateBitbucketRepoWithLfsOrFatal create, track, commit and push a random JPG file with git LFS, and return its hash.
func UpdateBitbucketRepoWithLfsOrFatal(t *testing.T, ctxt *pipelinectxt.ODSContext, wsDir, projectKey, filename string) [32]byte {
	fileHash := createJpgRandomFileOrFatal(t, wsDir, filename)
	trackLfsJpgFileToBitbucketOrFatal(t, wsDir, filename)
	// force reset commit SHA to latest commit that was just made
	ctxt.GitCommitSHA = ""
	assembleAndCacheOdsCtxtOrFatal(t, ctxt, wsDir)
	return fileHash
}

func createJpgRandomFileOrFatal(t *testing.T, wsDir, filename string) [32]byte {
	fileContent := make([]byte, 1000000) // 1MB
	_, err := rand.Read(fileContent)
	if err != nil {
		t.Fatalf("error creating random JPG file: %s", err)
	}
	err = writeFile(filepath.Join(wsDir, filename), string(fileContent))
	if err != nil {
		t.Fatalf("could not write file=%s: %s", filename, err)
	}
	return sha256.Sum256(fileContent)
}

func trackLfsJpgFileToBitbucketOrFatal(t *testing.T, wsDir, filename string) {
	stdout, stderr, err := command.RunInDir("git", []string{"lfs", "track", "*.jpg"}, wsDir)
	if err != nil {
		t.Fatalf("failed to track %s as LFS file: %s, stdout: %s, stderr: %s", filename, err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"add", ".gitattributes", filename}, wsDir)
	if err != nil {
		t.Fatalf("error running git add: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"commit", "-m", "track as LFS file " + filename}, wsDir)
	if err != nil {
		t.Fatalf("error running git commit: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	stdout, stderr, err = command.RunInDir("git", []string{"push", "origin", "master"}, wsDir)
	if err != nil {
		t.Fatalf("failed to push to remote: %s, stdout: %s, stderr: %s", err, stdout, stderr)
	}
}

func writeFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}
