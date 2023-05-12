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
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

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

func TestFetchODSConfig(t *testing.T) {
	bitbucketClient := &bitbucket.TestClient{}

	tests := map[string]struct {
		files   map[string][]byte
		wantErr string
		wantODS *config.ODS
	}{
		"no ODS file": {
			files:   map[string][]byte{},
			wantErr: "ods.yml not found",
		},
		"empty ODS file": {
			files: map[string][]byte{
				"ods.yaml": []byte(""),
			},
			wantErr: "config is empty",
		},
		"ods.yaml file": {
			files: map[string][]byte{
				"ods.yaml": []byte("pipelines: []"),
			},
			wantErr: "",
			wantODS: &config.ODS{Pipelines: []config.Pipeline{}},
		},
		"ods.yml file": {
			files: map[string][]byte{
				"ods.yml": []byte("pipelines: []"),
			},
			wantErr: "",
			wantODS: &config.ODS{Pipelines: []config.Pipeline{}},
		},
		"ods.yaml has precedence over ods.yml file": {
			files: map[string][]byte{
				"ods.yaml": []byte("pipelines: [{tasks: [{name: yaml}]}]"),
				"ods.yml":  []byte("pipelines: [{tasks: [{name: yml}]}]"),
			},
			wantErr: "",
			wantODS: &config.ODS{Pipelines: []config.Pipeline{{Tasks: []v1beta1.PipelineTask{{Name: "yaml"}}}}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bitbucketClient.Files = tc.files
			// As context is read from fake Bitbucket, dir value is unused.
			got, err := fetchODSConfig(bitbucketClient, "foo", "bar", "refs/heads/master")
			if tc.wantErr == "" && err != nil {
				t.Fatal(err)
			} else if tc.wantErr != "" && err == nil {
				t.Fatalf("want err: %s, got nothing", tc.wantErr)
			} else if err != nil && tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("want err: %s, got err: %s", tc.wantErr, err)
			}
			if diff := cmp.Diff(tc.wantODS, got); diff != "" {
				t.Fatalf("context mismatch (-want +got):\n%s", diff)
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
		"repo:refs_changed (branch push) without matching trigger is no-op": {
			requestBodyFixture: "manager/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/non-matching-trigger-ods.yaml"),
				},
			},
			wantStatus:         http.StatusAccepted,
			wantBody:           `{"title":"Accepted","detail":"Could not identify any pipeline to run as no trigger matched"}`,
			wantPipelineConfig: false,
		},
		"repo:refs_changed (branch push) triggers pipeline": {
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
		"repo:refs_changed (tag push) triggers pipeline": {
			requestBodyFixture: "manager/payload-tag.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/ods.yaml"),
				},
			},
			wantBody:           string(readTestdataFile(t, "golden/manager/response-payload-refs-changed-tag.json")),
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

func TestIdentifyPipelineConfig(t *testing.T) {
	tests := map[string]struct {
		pInfo     PipelineInfo
		odsConfig config.ODS
		// Index of pipeline that should be selected. -1 indicates no pipeline should be selected.
		wantPipelineIndex int
		// Index of trigger within pipeline that should be selected. -1 indicates no trigger should be selected.
		wantTriggerIndex int
	}{
		// branch
		"branch push - pipeline with no triggers": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  -1,
		},
		"branch push - pipeline with multiple triggers with no branch constraints": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{},
						{},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		"branch push - pipeline with one trigger with matching branch constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Branches: []string{"develop"}},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		"branch push - pipeline with one trigger with non-matching branch constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{{Branches: []string{"master"}}}},
				},
			},
			wantPipelineIndex: -1,
			wantTriggerIndex:  -1,
		},
		"branch push - pipeline with one trigger with tag constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{{Tags: []string{"v*"}}}},
				},
			},
			wantPipelineIndex: -1,
			wantTriggerIndex:  -1,
		},
		"branch push - pipeline with multiple triggers, one of them with matching branch constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Branches: []string{"master"}},
						{Branches: []string{"develop"}},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  1,
		},
		"branch push - pipeline with multiple triggers, none of them with matching branch constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Branches: []string{"master"}},
						{Branches: []string{"production"}},
					}},
				},
			},
			wantPipelineIndex: -1,
			wantTriggerIndex:  -1,
		},
		"branch push - pipeline with multiple triggers, all of them with matching branch constraints": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Branches: []string{"*"}},
						{Branches: []string{"develop"}},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		"branch push - multiple pipelines, one with matching trigger": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Branches: []string{"master"}},
					}},
					{Triggers: []config.Trigger{
						{Branches: []string{"*"}},
					}},
				},
			},
			wantPipelineIndex: 1,
			wantTriggerIndex:  0,
		},
		// tag
		"tag push - pipeline with multiple triggers with no tag constraints": {
			pInfo: PipelineInfo{ChangeRefType: "TAG", GitRef: "v1.0.0"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{},
						{},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		"tag push - pipeline with one trigger with matching tag constraint": {
			pInfo: PipelineInfo{ChangeRefType: "TAG", GitRef: "v1.0.0"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Tags: []string{"v*"}},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		"tag push - pipeline with one trigger with non-matching tag constraint": {
			pInfo: PipelineInfo{ChangeRefType: "TAG", GitRef: "v1.0.0"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Tags: []string{"v2.0.0"}},
					}},
				},
			},
			wantPipelineIndex: -1,
			wantTriggerIndex:  -1,
		},
		"tag push - pipeline with one trigger with branch constraint": {
			pInfo: PipelineInfo{ChangeRefType: "TAG", GitRef: "v1.0.0"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Branches: []string{"*"}},
					}},
				},
			},
			wantPipelineIndex: -1,
			wantTriggerIndex:  -1,
		},
		"tag push - pipeline with multiple triggers, one of them with matching tag constraint": {
			pInfo: PipelineInfo{ChangeRefType: "TAG", GitRef: "v1.0.0"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Tags: []string{"v2.0.0"}},
						{Tags: []string{"v1.0.0"}},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  1,
		},
		"tag push - pipeline with multiple triggers, none of them with matching tag constraint": {
			pInfo: PipelineInfo{ChangeRefType: "TAG", GitRef: "v1.0.0"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Tags: []string{"v2.0.0"}},
						{Tags: []string{"v3.0.0"}},
					}},
				},
			},
			wantPipelineIndex: -1,
			wantTriggerIndex:  -1,
		},
		"tag push - pipeline with multiple triggers, all of them with matching tag constraints": {
			pInfo: PipelineInfo{ChangeRefType: "TAG", GitRef: "v1.0.0"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Tags: []string{"v1.0.0"}},
						{Tags: []string{"*"}},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		"tag push - multiple pipelines, one with matching trigger": {
			pInfo: PipelineInfo{ChangeRefType: "TAG", GitRef: "v1.0.0"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Branches: []string{"master"}},
					}},
					{Triggers: []config.Trigger{
						{Tags: []string{"*"}},
					}},
				},
			},
			wantPipelineIndex: 1,
			wantTriggerIndex:  0,
		},
		// event
		"event - pipeline with non-matching event constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop", TriggerEvent: "pr:opened"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Events: []string{"repo:refs_pushed"}},
					}},
				},
			},
			wantPipelineIndex: -1,
			wantTriggerIndex:  -1,
		},
		"event - pipeline with matching event constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop", TriggerEvent: "pr:opened"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{Events: []string{"pr:opened"}},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		"event - pipeline with no event constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop", TriggerEvent: "pr:opened"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		// PR comment
		"PR comment - pipeline with non-matching PR comment constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop", Comment: "/deploy"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{PrComment: pString("/retest")},
					}},
				},
			},
			wantPipelineIndex: -1,
			wantTriggerIndex:  -1,
		},
		"PR comment - pipeline with matching PR comment constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop", Comment: "/deploy"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{PrComment: pString("/deploy")},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
		"PR comment - pipeline with no PR comment constraint": {
			pInfo: PipelineInfo{ChangeRefType: "BRANCH", GitRef: "develop", Comment: "/deploy"},
			odsConfig: config.ODS{
				Pipelines: []config.Pipeline{
					{Triggers: []config.Trigger{
						{},
					}},
				},
			},
			wantPipelineIndex: 0,
			wantTriggerIndex:  0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Annotate wanted pipeline/trigger so that
			// we can check later if it was selected.
			if tc.wantPipelineIndex > -1 {
				tc.odsConfig.Pipelines[tc.wantPipelineIndex].Tasks = []v1beta1.PipelineTask{
					{Name: "match this"},
				}
				if tc.wantTriggerIndex > -1 {
					tc.odsConfig.Pipelines[tc.wantPipelineIndex].Triggers[tc.wantTriggerIndex].Params = []v1beta1.Param{
						tektonStringParam("match", "this"),
					}
				}
			}
			got := identifyPipelineConfig(tc.pInfo, tc.odsConfig, "component")
			if tc.wantPipelineIndex > -1 && got == nil {
				t.Fatal("wanted a matching pipeline but got none")
			}
			if tc.wantPipelineIndex < 0 && got != nil {
				t.Fatal("wanted no matching pipeline, but got one")
			}
			if tc.wantPipelineIndex < 0 && got == nil {
				return // no matching pipeline, as wanted by the test case.
			}
			if len(got.PipelineSpec.Tasks) < 1 {
				t.Fatal("did not match wanted pipeline")
			}
			if tc.wantTriggerIndex > -1 && len(got.Params) < 1 {
				t.Fatal("did not match wanted trigger")
			}
		})
	}
}

func TestAnyPatternMatches(t *testing.T) {
	match := true
	tests := []struct {
		input   string
		pattern []string
		want    bool
	}{
		{"master", []string{"*"}, match},
		// TODO: The following is probably expected to work by users but does not work right now.
		// {"feature/foo", []string{"*"}, true},
		{"master", []string{"main", "*"}, match},
		{"feature/foo", []string{"feature/*"}, match},
		{"feature/foo", []string{"*/*"}, match},
		{"production", []string{}, match},
		{"feature/foo", []string{"feature"}, !match},
		{"production", []string{"main", "develop"}, !match},
		{"production", []string{"*/*"}, !match},
		{"production", []string{"p"}, !match},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v matches %s", tc.pattern, tc.input), func(t *testing.T) {
			got := anyPatternMatches(tc.input, tc.pattern)
			if got != tc.want {
				t.Fatalf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func pString(v string) *string { return &v }

func extractTaskNames(pConfig PipelineConfig) []string {
	gotNames := make([]string, 0, len(pConfig.PipelineSpec.Tasks))
	for _, task := range pConfig.PipelineSpec.Tasks {
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
