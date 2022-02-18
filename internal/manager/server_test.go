package manager

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/projectpath"
	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/internal/testfile"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const testWebhookSecret = "s3cr3t"

func TestRenderPipeline(t *testing.T) {
	wantPipeline := testfile.ReadGolden(t, "manager/pipeline.yaml")
	data := PipelineData{
		Name:            "bar-main",
		Project:         "foo",
		Repository:      "foo-bar",
		Component:       "bar",
		GitRef:          "main",
		GitFullRef:      "refs/heads/main",
		GitSHA:          "ef8755f06ee4b28c96a847a95cb8ec8ed6ddd1ca",
		RepoBase:        "https://bitbucket.acme.org",
		GitURI:          "https://bitbucket.acme.org/scm/foo/bar.git",
		Namespace:       "foo-cd",
		Stage:           "dev",
		TriggerEvent:    "repo:refs_changed",
		Comment:         "",
		PullRequestKey:  0,
		PullRequestBase: "",
	}

	// read ods.yaml
	conf := testfile.ReadFixture(t, "manager/ods.yaml")
	var odsConfig *config.ODS
	err := yaml.Unmarshal(conf, &odsConfig)
	fatalIfErr(t, err)
	gotPipeline, err := renderPipeline(odsConfig, data, "ClusterTask", "-v0-1-0")
	fatalIfErr(t, err)
	if diff := cmp.Diff(wantPipeline, gotPipeline); diff != "" {
		t.Fatalf("renderPipeline() mismatch (-want +got):\n%s", diff)
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

func TestShortenString(t *testing.T) {
	tests := map[string]struct {
		s        string
		max      int
		expected string
	}{
		"short enough": {
			s:        "foobar",
			max:      10,
			expected: "foobar",
		},
		"too long": {
			s:        "some-arbitarily-long-name-that-should-be-way-shorter",
			max:      30,
			expected: "some-arbitarily-long-n-8b85b7c",
		},
		"too long with slight difference in cut off string": {
			s:        "some-arbitarily-long-name-that-should-be-way-shorterx",
			max:      30,
			expected: "some-arbitarily-long-n-50a3b84",
		},
		"exact length": {
			s:        "some-arbitarily-long-name-that",
			max:      30,
			expected: "some-arbitarily-long-name-that",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := fitStringToMaxLength(tc.s, tc.max)
			if tc.expected != got {
				t.Fatalf(
					"Want '%s', got '%s' for (s='%s', max='%d')",
					tc.expected,
					got,
					tc.s,
					tc.max,
				)
			}
		})
	}
}

func TestMakePipelineName(t *testing.T) {
	tests := map[string]struct {
		component string
		branch    string
		expected  string
	}{
		"branch contains non-alphanumeric characters": {
			component: "comp",
			branch:    "bugfix/prj-529-bar-6-baz",
			expected:  "comp-bugfix-prj-529-bar-6-baz",
		},
		"branch contains uppercase characters": {
			component: "comp",
			branch:    "PRJ-529-bar-6-baz",
			expected:  "comp-prj-529-bar-6-baz",
		},
		"branch name is too long": {
			component: "comp",
			branch:    "bugfix/some-arbitarily-long-branch-name-that-should-be-way-shorter",
			expected:  "comp-bugfix-some-arbitarily-long-branch-name-th-87136df",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := makePipelineName(tc.component, tc.branch)
			if tc.expected != got {
				t.Fatalf(
					"Want '%s', got '%s' for (component='%s', branch='%s')",
					tc.expected,
					got,
					tc.component,
					tc.branch,
				)
			}
		})
	}
}

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

func fatalIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

type fakePruner struct {
	called chan bool
}

func (p *fakePruner) Prune(ctxt context.Context, pipelineRuns []tekton.PipelineRun) error {
	p.called <- true
	return nil
}

func testServer(kc kubernetes.ClientInterface, tc tektonClient.ClientInterface, bc bitbucketInterface, pruner PipelineRunPruner) (*httptest.Server, error) {
	server, err := NewServer(ServerConfig{
		Namespace:     "bar-cd",
		Project:       "bar",
		Token:         "test",
		WebhookSecret: testWebhookSecret,
		TaskKind:      "ClusterTask",
		RepoBase:      "https://domain.com",
		StorageConfig: StorageConfig{
			Provisioner: "kubernetes.io/aws-ebs",
			ClassName:   "gp2",
			Size:        "2Gi",
		},
		KubernetesClient:  kc,
		TektonClient:      tc,
		BitbucketClient:   bc,
		PipelineRunPruner: pruner,
	})
	if err != nil {
		return nil, err
	}
	return httptest.NewServer(http.HandlerFunc(server.HandleRoot)), nil
}

func TestWebhookHandling(t *testing.T) {

	tests := map[string]struct {
		requestBodyFixture string
		kubernetesClient   *kubernetes.TestClient
		tektonClient       *tektonClient.TestClient
		bitbucketClient    *bitbucket.TestClient
		wrongSignature     bool
		wantStatus         int
		wantBody           string
		check              func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner)
	}{
		"wrong signature is not processed": {
			requestBodyFixture: "manager/payload.json", // valid payload
			wrongSignature:     true,
			wantStatus:         http.StatusBadRequest,
			wantBody:           "failed to validate incoming request",
			check: func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
		},
		"invalid JSON is not processed": {
			requestBodyFixture: "manager/payload-invalid.json",
			wantStatus:         http.StatusBadRequest,
			wantBody:           "cannot parse JSON: invalid character '\\n' in string literal",
			check: func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
		},
		"unsupported events are not processed": {
			requestBodyFixture: "manager/payload-unknown-event.json",
			wantStatus:         http.StatusBadRequest,
			wantBody:           "Unsupported event key: repo:ref_changed",
			check: func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
		},
		"tags are not processed": {
			requestBodyFixture: "manager/payload-tag.json",
			wantStatus:         http.StatusTeapot,
			wantBody:           "Skipping change ref type TAG, only BRANCH is supported",
			check: func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
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
			wantStatus: http.StatusTeapot,
			wantBody:   "Commit should be skipped",
			check: func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
		},
		"pushes into new branch creates a pipeline": {
			requestBodyFixture: "manager/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/ods.yaml"),
				},
			},
			wantBody:   string(readTestdataFile(t, "golden/manager/response-payload-refs-changed.json")),
			wantStatus: http.StatusOK,
			check: func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner) {
				if len(tc.CreatedPipelines) != 1 || len(tc.UpdatedPipelines) != 0 {
					t.Fatal("exactly one pipeline should have been created")
				}
				if len(kc.CreatedPVCs) != 1 {
					t.Fatal("exactly one PVC should have been created")
				}
				if len(tc.CreatedPipelineRuns) != 1 {
					t.Fatal("exactly one pipeline run should have been created")
				}
				select {
				case <-p.called:
					t.Log("pruning occured")
				case <-time.After(10 * time.Second):
					t.Fatal("pruning should have occured")
				}
			},
		},
		"pushes into an existing branch updates a pipeline": {
			requestBodyFixture: "manager/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/ods.yaml"),
				},
			},
			kubernetesClient: &kubernetes.TestClient{
				PVCs: []*corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							// generated PVC name
							Name: "ods-workspace-bar",
						},
					},
				},
			},
			tektonClient: &tektonClient.TestClient{
				Pipelines: []*tekton.Pipeline{
					{
						ObjectMeta: metav1.ObjectMeta{
							// generated pipeline name
							Name: "bar-master",
						},
					},
				},
			},
			wantBody:   string(readTestdataFile(t, "golden/manager/response-payload-refs-changed.json")),
			wantStatus: http.StatusOK,
			check: func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner) {
				if len(tc.CreatedPipelines) != 0 || len(tc.UpdatedPipelines) != 1 {
					t.Fatal("exactly one pipeline should have been updated")
				}
				if len(kc.CreatedPVCs) > 0 {
					t.Fatal("no PVC should have been created")
				}
				if len(tc.CreatedPipelineRuns) != 1 {
					t.Fatal("exactly one pipeline run should have been created")
				}
				select {
				case <-p.called:
					t.Log("pruning occured")
				case <-time.After(10 * time.Second):
					t.Fatal("pruning should have occured")
				}
			},
		},
		"PR open events update a pipeline": {
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
			tektonClient: &tektonClient.TestClient{
				Pipelines: []*tekton.Pipeline{
					{
						ObjectMeta: metav1.ObjectMeta{
							// generated pipeline name
							Name: "bar-feature-foo",
						},
					},
				},
			},
			wantBody:   string(readTestdataFile(t, "golden/manager/response-payload-pr-opened.json")),
			wantStatus: http.StatusOK,
			check: func(t *testing.T, kc *kubernetes.TestClient, tc *tektonClient.TestClient, bc *bitbucket.TestClient, p *fakePruner) {
				if len(tc.CreatedPipelines) != 0 || len(tc.UpdatedPipelines) != 1 {
					t.Fatal("exactly one pipeline should have been updated")
				}
				if len(tc.CreatedPipelineRuns) != 1 {
					t.Fatal("exactly one pipeline run should have been created")
				}
				select {
				case <-p.called:
					t.Log("pruning occured")
				case <-time.After(10 * time.Second):
					t.Fatal("pruning should have occured")
				}
			},
		},
		"failure to create pipeline is handled properly": {
			requestBodyFixture: "manager/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/ods.yaml"),
				},
			},
			tektonClient: &tektonClient.TestClient{
				FailCreatePipeline: true,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "cannot create pipeline bar-master",
		},
		"failure to update pipeline is handled properly": {
			requestBodyFixture: "manager/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/manager/ods.yaml"),
				},
			},
			tektonClient: &tektonClient.TestClient{
				Pipelines: []*tekton.Pipeline{
					{
						ObjectMeta: metav1.ObjectMeta{
							// generated pipeline name
							Name: "bar-master",
						},
					},
				},
				FailUpdatePipeline: true,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "cannot update pipeline bar-master",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.tektonClient == nil {
				tc.tektonClient = &tektonClient.TestClient{}
			}
			if tc.bitbucketClient == nil {
				tc.bitbucketClient = &bitbucket.TestClient{}
			}
			if tc.kubernetesClient == nil {
				tc.kubernetesClient = &kubernetes.TestClient{}
			}
			pruner := &fakePruner{called: make(chan bool)}
			ts, err := testServer(tc.kubernetesClient, tc.tektonClient, tc.bitbucketClient, pruner)
			if err != nil {
				t.Fatal(err)
			}
			defer ts.Close()
			filename := filepath.Join(projectpath.Root, "test/testdata/fixtures", tc.requestBodyFixture)
			f, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}
			body, err := ioutil.ReadAll(f)
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
			gotBodyBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			gotBody := removeSpace(string(gotBodyBytes))
			if diff := cmp.Diff(removeSpace(tc.wantBody), gotBody); diff != "" {
				t.Fatalf("body mismatch (-want +got):\n%s", diff)
			}
			if tc.check != nil {
				tc.check(t, tc.kubernetesClient, tc.tektonClient, tc.bitbucketClient, pruner)
			}
		})
	}
}

func renderPipeline(odsConfig *config.ODS, data PipelineData, taskKind tekton.TaskKind, taskSuffix string) ([]byte, error) {
	p := assemblePipeline(odsConfig, data, taskKind, taskSuffix)
	return yaml.Marshal(p)
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
	b, err := ioutil.ReadFile(filepath.Join(projectpath.Root, "test/testdata", filename))
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
