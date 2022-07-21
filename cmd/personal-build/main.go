package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/manager"
	"github.com/opendevstack/pipeline/pkg/logging"
)

const signatureHeader = "X-Hub-Signature"

type options struct {
	// Location of secrets file of Helm chart of the project.
	secretsFile string
	project     string
	url         string
	debug       bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.secretsFile, "secrets-file", "", "Secrets file in project helm chart dir. Used to extract bitbucketWebhookSecret")
	flag.StringVar(&opts.project, "project", "", "Project name")
	flag.StringVar(&opts.url, "url", "", "URL of pipeline manager")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "debug mode")
	flag.Parse()

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	stdout, stderr, err := command.Run("git", []string{"symbolic-ref", "--short", "HEAD"})
	if err != nil {
		log.Fatal(fmt.Errorf("cannot get branch name: %s (%w)", stderr, err))
	}
	branchName := strings.TrimSpace(string(stdout))
	logger.Infof("branchName=%s", branchName)

	stdout, stderr, err = command.Run("git", []string{"symbolic-ref", "HEAD"})
	if err != nil {
		log.Fatal(fmt.Errorf("cannot get full ref name: %s (%w)", stderr, err))
	}
	fullRefName := strings.TrimSpace(string(stdout))
	logger.Infof("fullRefName=%s", fullRefName)

	stdout, stderr, err = command.Run("git", []string{"rev-parse", "HEAD"})
	if err != nil {
		log.Fatal(fmt.Errorf("cannot get commit sha: %s (%w)", stderr, err))
	}
	toHash := strings.TrimSpace(string(stdout))
	logger.Infof("toHash=%s", toHash)

	stdout, stderr, err = command.Run("sops", []string{"-d", "--extract", `["setup"]["bitbucketWebhookSecret"]`, opts.secretsFile})
	if err != nil {
		log.Fatal(fmt.Errorf("cannot get bitbucketWebhookSecret: %s (%w)", stderr, err))
	}
	webhookSecret := strings.TrimSpace(string(stdout))
	logger.Infof("Read bitbucketWebhookSecret without errors")
	logger.Debugf("bitbucketWebhookSecret=%s", webhookSecret)

	var path string
	if path, err = os.Getwd(); err != nil {
		log.Fatal(fmt.Errorf("Cannot get current work directory: %s (%w)", stderr, err))
	}
	dirname := strings.ToLower(filepath.Base(path))
	prefix := strings.ToLower(opts.project) + "-"
	if !strings.HasPrefix(dirname, prefix) {
		log.Fatal(fmt.Errorf("current working directory must have project prefix: %s (does not have prefix %s)", dirname, prefix))
	}
	req := &manager.RequestPersonalBuild{
		RepositorySlug: dirname,
		FullRefName:    fullRefName,
		BranchName:     branchName,
		ToHash:         toHash,
	}

	postPersonalBuild(logger, opts.url, req, webhookSecret)
}

func postPersonalBuild(logger logging.LeveledLoggerInterface, url string, personalBuildRequest *manager.RequestPersonalBuild, secret string) {
	postBody, err := json.Marshal(personalBuildRequest)
	if err != nil {
		log.Fatal(fmt.Errorf("Could not construct json body for post request: %w", err))
	}
	logger.Infof("Json post body: %s", string(postBody[:]))
	req, err := http.NewRequest("POST", url, bytes.NewReader(postBody))
	if err != nil {
		log.Fatalf("NewRequest: %v", err)
	}
	hmacValue := hmacHeader(secret, postBody)
	logger.Infof("%s=%s", signatureHeader, hmacValue)
	req.Header.Set(signatureHeader, hmacValue)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Minute}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Http request failed: %v", err)
	}
	logger.Infof("%s response status=%d", url, res.StatusCode)
	gotBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%s response body could not be read", url)
	}
	logger.Infof("%s response body=%s", url, string(gotBodyBytes))
}

// hmacHeader generates a X-Hub-Signature header given a secret token and the request body
// See https://developer.github.com/webhooks/securing/#validating-payloads-from-github
// Note that while this example and the validation comes from GitHub, it applies to
// Bitbucket just the same.
func hmacHeader(secret string, body []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write(body)
	if err != nil {
		log.Fatalf("HMACHeader fail: %s", err)
	}
	return fmt.Sprintf("sha256=%s", hex.EncodeToString(h.Sum(nil)))
}
