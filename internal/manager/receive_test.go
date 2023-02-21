package manager

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/httpjson"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
)

func TestSelectEnvironmentFromMapping(t *testing.T) {
	tests := []struct {
		mapping []config.BranchToEnvironmentMapping
		branch  string
		want    string
	}{
		{[]config.BranchToEnvironmentMapping{
			{
				Branch:      "develop",
				Environment: "dev",
			},
		}, "develop", "dev"},
		{[]config.BranchToEnvironmentMapping{
			{
				Branch:      "develop",
				Environment: "dev",
			},
		}, "developer", ""},
		{[]config.BranchToEnvironmentMapping{
			{
				Branch:      "develop",
				Environment: "dev",
			},
			{
				Branch:      "develop",
				Environment: "foo",
			},
		}, "develop", "dev"},
		{[]config.BranchToEnvironmentMapping{
			{
				Branch:      "release/*",
				Environment: "qa",
			},
		}, "release/1.0", "qa"},
		{[]config.BranchToEnvironmentMapping{
			{
				Branch:      "release/*",
				Environment: "qa",
			},
		}, "release", ""},
		{[]config.BranchToEnvironmentMapping{
			{
				Branch:      "*",
				Environment: "dev",
			},
		}, "foo", "dev"},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("mapping #%d", i), func(t *testing.T) {
			got := selectEnvironmentFromMapping(tc.mapping, tc.branch)
			if tc.want != got {
				t.Fatalf("Got %v, want %v for branch '%s'", got, tc.want, tc.branch)
			}
		})
	}
}

func TestIsCiSkipInCommitMessage(t *testing.T) {
	tests := []struct {
		message string
		want    bool
	}{
		{"docs: update README [ci skip]", true},
		{"docs: update README [skip ci]", true},
		{"docs: update README ***NO_CI***", true},
		{"docs: update READM", false},
		{"docs: update README\n\n- typo\n- [ci skip]", false},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("commit message #%d", i), func(t *testing.T) {
			got := isCiSkipInCommitMessage((tc.message))
			if tc.want != got {
				t.Fatalf("Got %v, want %v for message '%s'", got, tc.want, tc.message)
			}
		})
	}
}

func testServer(bc bitbucketInterface, ch chan PipelineConfig) *httptest.Server {
	r := &BitbucketWebhookReceiver{
		TriggeredPipelines: ch,
		Namespace:          "bar-cd",
		Project:            "bar",
		WebhookSecret:      testWebhookSecret,
		RepoBase:           "https://domain.com",
		BitbucketClient:    bc,
		Logger:             &logging.LeveledLogger{Level: logging.LevelNull},
	}
	return httptest.NewServer(httpjson.Handler(r.Handle))
}

func TestWebhookHandling(t *testing.T) {
	tests := map[string]struct {
		requestBodyFixture    string
		bitbucketClient       *bitbucket.TestClient
		wrongSignature        bool
		wantStatus            int
		wantBody              string
		wantPipelineConfig    bool
		wantPipelineTaskNames []string
	}{
		"wrong signature is not processed": {
			requestBodyFixture: "manager/payload.json", // valid payload
			wrongSignature:     true,
			wantStatus:         http.StatusUnauthorized,
			wantBody:           `{"title":"Unauthorized","detail":"failed to validate incoming request"}`,
			wantPipelineConfig: false,
		},
		"invalid JSON is not processed": {
			requestBodyFixture: "manager/payload-invalid.json",
			wantStatus:         http.StatusBadRequest,
			wantBody:           `{"title":"BadRequest","detail":"cannot parse JSON: invalid character '\\n' in string literal"}`,
			wantPipelineConfig: false,
		},
		"unsupported events are not processed": {
			requestBodyFixture: "manager/payload-unknown-event.json",
			wantStatus:         http.StatusBadRequest,
			wantBody:           `{"title":"BadRequest","detail":"unsupported event key: repo:ref_changed"}`,
			wantPipelineConfig: false,
		},
		"tags are not processed": {
			requestBodyFixture: "manager/payload-tag.json",
			wantStatus:         http.StatusUnprocessableEntity,
			wantBody:           `{"title":"UnprocessableEntity","detail":"skipping change ref type TAG, only BRANCH is supported"}`,
			wantPipelineConfig: false,
		},
		"commits with skip message are not processed": {
			requestBodyFixture: "manager/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Commits: []bitbucket.Commit{
					{
						// commit referenced in payload
						ID:      "0e183aa3bc3c6deb8f40b93fb2fc4354533cf62f",
						Message: "Update readme [ci skip]",
					},
				},
			},
			wantStatus:         http.StatusAccepted,
			wantBody:           `{"title":"Accepted","detail":"Commit 0e183aa3bc3c6deb8f40b93fb2fc4354533cf62f should be skipped"}`,
			wantPipelineConfig: false,
		},
		"repo:refs_changed triggers pipeline": {
			requestBodyFixture: "manager/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/ods.yaml"),
				},
			},
			wantBody:           string(readTestdataFile(t, "golden/manager/response-payload-refs-changed.json")),
			wantStatus:         http.StatusOK,
			wantPipelineConfig: true,
		},
		"pr:opened triggers pipeline": {
			requestBodyFixture: "manager/payload-pr-opened.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/ods.yaml"),
				},
				PullRequests: []bitbucket.PullRequest{
					{
						Open: true,
						ID:   1,
						ToRef: bitbucket.Ref{
							ID: "refs/heads/master",
						},
					},
				},
			},
			wantBody:           string(readTestdataFile(t, "golden/manager/response-payload-pr-opened.json")),
			wantStatus:         http.StatusOK,
			wantPipelineConfig: true,
		},
		"pr:comment:added request triggers pipeline with matching comment": {
			requestBodyFixture: "manager/payload-pr-comment-added-select.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/multi-pipeline-ods.yaml"),
				},
				PullRequests: []bitbucket.PullRequest{
					{
						Open: true,
						ID:   1,
						ToRef: bitbucket.Ref{
							ID: "refs/heads/master",
						},
					},
				},
			},
			wantBody:              string(readTestdataFile(t, "golden/manager/response-payload-pr-comment-added-select.json")),
			wantStatus:            http.StatusOK,
			wantPipelineConfig:    true,
			wantPipelineTaskNames: []string{"go-helm-build-comment-added-select-foo"},
		},
		"pr:comment:added request on excluded branch triggers branch-specific pipeline": {
			requestBodyFixture: "manager/payload-pr-comment-added-select-foo.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/multi-pipeline-ods.yaml"),
				},
				PullRequests: []bitbucket.PullRequest{
					{
						Open: true,
						ID:   1,
						ToRef: bitbucket.Ref{
							ID: "refs/heads/master",
						},
					},
				},
			},
			wantBody:              string(readTestdataFile(t, "golden/manager/response-payload-pr-comment-added-select-foo.json")),
			wantStatus:            http.StatusOK,
			wantPipelineConfig:    true,
			wantPipelineTaskNames: []string{"go-helm-build-opened-pr-foo"},
		},
		"pr:comment:added request triggers catch-all pipeline": {
			requestBodyFixture: "manager/payload-pr-comment-added-other.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/multi-pipeline-ods.yaml"),
				},
				PullRequests: []bitbucket.PullRequest{
					{
						Open: true,
						ID:   1,
						ToRef: bitbucket.Ref{
							ID: "refs/heads/master",
						},
					},
				},
			},
			wantBody:              string(readTestdataFile(t, "golden/manager/response-payload-pr-comment-added-other.json")),
			wantStatus:            http.StatusOK,
			wantPipelineConfig:    true,
			wantPipelineTaskNames: []string{"go-helm-build-catch-all"},
		},
		"pr:opened on feature/foo request triggers matching pipeline": {
			requestBodyFixture: "manager/payload-pr-opened.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/multi-pipeline-ods.yaml"),
				},
				PullRequests: []bitbucket.PullRequest{
					{
						Open: true,
						ID:   1,
						ToRef: bitbucket.Ref{
							ID: "refs/heads/master",
						},
					},
				},
			},
			wantBody:              string(readTestdataFile(t, "golden/manager/response-payload-pr-opened-feature-foo.json")),
			wantStatus:            http.StatusOK,
			wantPipelineConfig:    true,
			wantPipelineTaskNames: []string{"go-helm-build-opened-pr-foo"},
		},
		"pr:opened on feature/other request triggers matching pipeline": {
			requestBodyFixture: "manager/payload-pr-opened-feature-other.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/multi-pipeline-ods.yaml"),
				},
				PullRequests: []bitbucket.PullRequest{
					{
						Open: true,
						ID:   1,
						ToRef: bitbucket.Ref{
							ID: "refs/heads/master",
						},
					},
				},
			},
			wantBody:              string(readTestdataFile(t, "golden/manager/response-payload-pr-opened-feature-other.json")),
			wantStatus:            http.StatusOK,
			wantPipelineConfig:    true,
			wantPipelineTaskNames: []string{"go-helm-build-opened-pr"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.bitbucketClient == nil {
				tc.bitbucketClient = &bitbucket.TestClient{}
			}
			// Allow to send one PipelineConfig to the channel without blocking
			ch := make(chan PipelineConfig, 1)
			ts := testServer(tc.bitbucketClient, ch)
			defer ts.Close()
			filename := filepath.Join(projectpath.Root, "test/testdata/fixtures", tc.requestBodyFixture)
			f, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			fr := bytes.NewReader(body)
			req, err := http.NewRequest("POST", ts.URL, fr)
			if err != nil {
				t.Fatalf("NewRequest: %v", err)
			}
			if tc.wrongSignature {
				req.Header.Set(signatureHeader, "foobar")
			} else {
				req.Header.Set(signatureHeader, hmacHeader(t, testWebhookSecret, body))
			}
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{Timeout: time.Minute}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			gotStatus := res.StatusCode
			if tc.wantStatus != gotStatus {
				t.Fatalf("Got status: %v, want: %v", gotStatus, tc.wantStatus)
			}
			gotBodyBytes, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			gotBody := removeSpace(string(gotBodyBytes))
			if diff := cmp.Diff(removeSpace(tc.wantBody), gotBody); diff != "" {
				t.Fatalf("body mismatch (-want +got):\n%s", diff)
			}
			// Check if request sent a pipeline config to ch.
			select {
			case pConfig := <-ch:
				if !tc.wantPipelineConfig {
					t.Fatal("want no pipeline config, got one")
				}
				if tc.wantPipelineTaskNames != nil {
					gotNames := extractTaskNames(pConfig)
					if diff := cmp.Diff(tc.wantPipelineTaskNames, gotNames); diff != "" {
						t.Fatalf("pipeline config mismatch (-want, +got):\n%s", diff)
					}
				}
			default:
				if tc.wantPipelineConfig {
					t.Fatal("want pipeline config, got none")
				}
			}
		})
	}
}

func extractTaskNames(pConfig PipelineConfig) []string {
	gotNames := make([]string, 0, len(pConfig.Tasks))
	for _, task := range pConfig.Tasks {
		gotNames = append(gotNames, task.Name)
	}
	return gotNames
}

func removeSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func readTestdataFile(t *testing.T, filename string) []byte {
	b, err := os.ReadFile(filepath.Join(projectpath.Root, "test/testdata", filename))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// hmacHeader generates a X-Hub-Signature header given a secret token and the request body
// See https://developer.github.com/webhooks/securing/#validating-payloads-from-github
// Note that while this example and the validation comes from GitHub, it applies to
// Bitbucket just the same.
func hmacHeader(t *testing.T, secret string, body []byte) string {
	t.Helper()
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write(body)
	if err != nil {
		t.Fatalf("HMACHeader fail: %s", err)
	}
	return fmt.Sprintf("sha256=%s", hex.EncodeToString(h.Sum(nil)))
}
