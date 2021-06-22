package interceptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/pkg/config"
	"sigs.k8s.io/yaml"
)

func TestRenderPipeline(t *testing.T) {
	wantPipeline := ReadGoldenFile(t, "pipeline.yml")
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
	conf := ReadFixtureFile(t, "ods.yml")
	var odsConfig config.ODS
	err := yaml.Unmarshal(conf, &odsConfig)
	fatalIfErr(t, err)
	phases := odsConfig.Phases
	phasesList := []config.Phases{phases}
	gotPipeline, err := renderPipeline(phasesList, data)
	fatalIfErr(t, err)
	if diff := cmp.Diff(wantPipeline, gotPipeline); diff != "" {
		t.Errorf("renderPipeline() mismatch (-want +got):\n%s", diff)
	}
}

func TestExtensions(t *testing.T) {
	bodyFixture := ReadFixtureFile(t, "payload.json")
	wantBody := ReadGoldenFile(t, "payload.json")
	data := PipelineData{
		Name:            "bar-main",
		Namespace:       "foo-cd",
		Project:         "foo",
		Repository:      "foo-bar",
		Component:       "bar",
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
		t.Errorf("extendBodyWithExtensions() mismatch (-want +got):\n%s", diff)
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

func fatalIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// ReadFixtureFile returns the contents of the fixture file or fails.
func ReadFixtureFile(t *testing.T, filename string) []byte {
	return readFileOrFatal(t, "../../test/testdata/fixtures/interceptor/"+filename)
}

// ReadGoldenFile returns the contents of the golden file or fails.
func ReadGoldenFile(t *testing.T, filename string) []byte {
	return readFileOrFatal(t, "../../test/testdata/golden/interceptor/"+filename)
}

func readFileOrFatal(t *testing.T, name string) []byte {
	b, err := readFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func readFile(name string) ([]byte, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return []byte{}, fmt.Errorf("Could not get filename when looking for %s", name)
	}
	filepath := path.Join(path.Dir(filename), name)
	return ioutil.ReadFile(filepath)
}
