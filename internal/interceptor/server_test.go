package interceptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/projectpath"
	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/internal/testfile"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestRenderPipeline(t *testing.T) {
	wantPipeline := testfile.ReadGolden(t, "interceptor/pipeline.yaml")
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
		TriggerEvent:    "repo:refs_changed",
		Comment:         "",
		PullRequestKey:  0,
		PullRequestBase: "",
	}

	// read ods.yaml
	conf := testfile.ReadFixture(t, "interceptor/ods.yaml")
	var odsConfig *config.ODS
	err := yaml.Unmarshal(conf, &odsConfig)
	fatalIfErr(t, err)
	gotPipeline, err := renderPipeline(odsConfig, data, "ClusterTask", "-v0-1-0")
	fatalIfErr(t, err)
	if diff := cmp.Diff(wantPipeline, gotPipeline); diff != "" {
		t.Fatalf("renderPipeline() mismatch (-want +got):\n%s", diff)
	}
}

func TestExtensions(t *testing.T) {
	bodyFixture := testfile.ReadFixture(t, "interceptor/payload.json")
	wantBody := testfile.ReadGolden(t, "interceptor/payload.json")
	data := PipelineData{
		Name:            "bar-main",
		Namespace:       "foo-cd",
		Project:         "foo",
		Repository:      "foo-bar",
		Component:       "bar",
		Environment:     "",
		Version:         "",
		GitRef:          "main",
		GitFullRef:      "refs/heads/main",
		GitSHA:          "ef8755f06ee4b28c96a847a95cb8ec8ed6ddd1ca",
		RepoBase:        "https://bitbucket.acme.org",
		GitURI:          "https://bitbucket.acme.org/scm/foo/bar.git",
		PVC:             "pipeline-bar",
		TriggerEvent:    "repo:refs_changed",
		Comment:         "",
		PullRequestKey:  0,
		PullRequestBase: "",
	}
	gotBody, err := extendBodyWithExtensions(bodyFixture, data)
	fatalIfErr(t, err)
	var wantPayload map[string]interface{}
	err = json.Unmarshal(wantBody, &wantPayload)
	fatalIfErr(t, err)
	var gotPayload map[string]interface{}
	err = json.Unmarshal(gotBody, &gotPayload)
	fatalIfErr(t, err)
	if diff := cmp.Diff(wantPayload, gotPayload); diff != "" {
		t.Fatalf("extendBodyWithExtensions() mismatch (-want +got):\n%s", diff)
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

func testServer(tc tektonClient.ClientInterface, bc bitbucketInterface) (*httptest.Server, error) {
	server, err := NewServer(ServerConfig{
		Namespace:       "bar-cd",
		Project:         "bar",
		Token:           "test",
		TaskKind:        "ClusterTask",
		RepoBase:        "https://domain.com",
		TektonClient:    tc,
		BitbucketClient: bc,
	})
	if err != nil {
		return nil, err
	}
	return httptest.NewServer(http.HandlerFunc(server.HandleRoot)), nil
}

func TestWebhookHandling(t *testing.T) {

	tests := map[string]struct {
		requestBodyFixture string
		tektonClient       *tektonClient.TestClient
		bitbucketClient    *bitbucket.TestClient
		wantStatus         int
		wantBody           string
		check              func(t *testing.T, tc *tektonClient.TestClient, bc *bitbucket.TestClient)
	}{
		"invalid JSON is not processed": {
			requestBodyFixture: "interceptor/payload-invalid.json",
			wantStatus:         http.StatusBadRequest,
			wantBody:           "cannot parse JSON: invalid character '\\n' in string literal",
			check: func(t *testing.T, tc *tektonClient.TestClient, bc *bitbucket.TestClient) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
		},
		"unsupported events are not processed": {
			requestBodyFixture: "interceptor/payload-unknown-event.json",
			wantStatus:         http.StatusBadRequest,
			wantBody:           "Unsupported event key: repo:ref_changed",
			check: func(t *testing.T, tc *tektonClient.TestClient, bc *bitbucket.TestClient) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
		},
		"tags are not processed": {
			requestBodyFixture: "interceptor/payload-tag.json",
			wantStatus:         http.StatusTeapot,
			wantBody:           "Skipping change ref type TAG, only BRANCH is supported",
			check: func(t *testing.T, tc *tektonClient.TestClient, bc *bitbucket.TestClient) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
		},
		"commits with skip message are not processed": {
			requestBodyFixture: "interceptor/payload.json",
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
			check: func(t *testing.T, tc *tektonClient.TestClient, bc *bitbucket.TestClient) {
				if len(tc.CreatedPipelines) > 0 || len(tc.UpdatedPipelines) > 0 {
					t.Fatal("no pipeline should have been created/updated")
				}
			},
		},
		"pushes into new branch creates a pipeline": {
			requestBodyFixture: "interceptor/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/interceptor/ods.yaml"),
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   string(readTestdataFile(t, "golden/interceptor/extended-payload.json")),
			check: func(t *testing.T, tc *tektonClient.TestClient, bc *bitbucket.TestClient) {
				if len(tc.CreatedPipelines) != 1 || len(tc.UpdatedPipelines) != 0 {
					t.Fatal("exactly one pipeline should have been created")
				}
			},
		},
		"pushes into an existing branch updates a pipeline": {
			requestBodyFixture: "interceptor/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/interceptor/ods.yaml"),
				},
			},
			tektonClient: &tektonClient.TestClient{
				Pipelines: []*tekton.Pipeline{
					{
						ObjectMeta: metav1.ObjectMeta{
							// generated pipeline name
							Name: "ods-pipeline-master",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   string(readTestdataFile(t, "golden/interceptor/extended-payload.json")),
			check: func(t *testing.T, tc *tektonClient.TestClient, bc *bitbucket.TestClient) {
				if len(tc.CreatedPipelines) != 0 || len(tc.UpdatedPipelines) != 1 {
					t.Fatal("exactly one pipeline should have been updated")
				}
			},
		},
		"PR open events update a pipeline": {
			requestBodyFixture: "interceptor/payload-pr-opened.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/interceptor/ods.yaml"),
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
							Name: "ods-pipeline-feature-foo",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   string(readTestdataFile(t, "golden/interceptor/extended-payload-pr-opened.json")),
			check: func(t *testing.T, tc *tektonClient.TestClient, bc *bitbucket.TestClient) {
				if len(tc.CreatedPipelines) != 0 || len(tc.UpdatedPipelines) != 1 {
					t.Fatal("exactly one pipeline should have been updated")
				}
			},
		},
		"failure to create pipeline is handled properly": {
			requestBodyFixture: "interceptor/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/interceptor/ods.yaml"),
				},
			},
			tektonClient: &tektonClient.TestClient{
				FailCreatePipeline: true,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "cannot create pipeline ods-pipeline-master",
		},
		"failure to update pipeline is handled properly": {
			requestBodyFixture: "interceptor/payload.json",
			bitbucketClient: &bitbucket.TestClient{
				Files: map[string][]byte{
					"ods.yaml": readTestdataFile(t, "fixtures/interceptor/ods.yaml"),
				},
			},
			tektonClient: &tektonClient.TestClient{
				Pipelines: []*tekton.Pipeline{
					{
						ObjectMeta: metav1.ObjectMeta{
							// generated pipeline name
							Name: "ods-pipeline-master",
						},
					},
				},
				FailUpdatePipeline: true,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "cannot update pipeline ods-pipeline-master",
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
			ts, err := testServer(tc.tektonClient, tc.bitbucketClient)
			if err != nil {
				t.Fatal(err)
			}
			defer ts.Close()
			filename := filepath.Join(projectpath.Root, "test/testdata/fixtures", tc.requestBodyFixture)
			f, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}
			res, err := http.Post(ts.URL, "application/json", f)
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
				tc.check(t, tc.tektonClient, tc.bitbucketClient)
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
