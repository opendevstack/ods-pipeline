package interceptor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/testfile"
	"github.com/opendevstack/pipeline/pkg/config"
	"sigs.k8s.io/yaml"
)

func TestRenderPipeline(t *testing.T) {
	wantPipeline := testfile.ReadGolden(t, "interceptor/pipeline.yml")
	data := PipelineData{
		Name:            "bar-main",
		Project:         "foo",
		Repository:      "foo-bar",
		Component:       "bar",
		GitRef:          "main",
		GitFullRef:      "refs/heads/main",
		GitSHA:          "ef8755f06ee4b28c96a847a95cb8ec8ed6ddd1ca",
		ResourceVersion: 0,
		RepoBase:        "https://bitbucket.acme.org",
		GitURI:          "https://bitbucket.acme.org/scm/foo/bar.git",
		Namespace:       "foo-cd",
		TriggerEvent:    "repo:refs_changed",
		Comment:         "",
		PullRequestKey:  0,
		PullRequestBase: "",
	}

	// read ods.yml
	conf := testfile.ReadFixture(t, "interceptor/ods.yml")
	var odsConfig *config.ODS
	err := yaml.Unmarshal(conf, &odsConfig)
	fatalIfErr(t, err)
	gotPipeline, err := renderPipeline(odsConfig, data)
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
		ResourceVersion: 0,
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
