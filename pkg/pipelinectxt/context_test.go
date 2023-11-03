package pipelinectxt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/ods-pipeline/internal/gittest"
)

func TestAssemble(t *testing.T) {
	c := &ODSContext{
		Namespace: "foo-cd",
	}

	dir, cleanup, err := gittest.CreateFakeGitRepoDir(
		"https://example.bitbucket.com/scm/ODS/ods-pipeline.git",
		"master",
		"7f96ec9fcf097e5b21687d402bc70370ac247d8a",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	err = c.Assemble(dir, "foo-n8j4k")
	if err != nil {
		t.Fatal(err)
	}
	wantContext := &ODSContext{
		Namespace:       "foo-cd",
		Project:         "ods",
		Repository:      "ods-pipeline",
		Component:       "pipeline",
		GitCommitSHA:    "7f96ec9fcf097e5b21687d402bc70370ac247d8a",
		GitFullRef:      "refs/heads/master",
		GitRef:          "master",
		GitURL:          "https://example.bitbucket.com/scm/ODS/ods-pipeline.git",
		PullRequestBase: "",
		PullRequestKey:  "",
		PipelineRun:     "foo-n8j4k",
	}
	if diff := cmp.Diff(wantContext, c); diff != "" {
		t.Fatalf("context mismatch (-want +got):\n%s", diff)
	}
}

func TestCopy(t *testing.T) {
	c1 := &ODSContext{
		Namespace:    "foo-cd",
		GitFullRef:   "refs/heads/master",
		GitURL:       "https://github.com/opendevstack/ods-pipeline.git",
		GitCommitSHA: "abcdef",
	}
	c2 := c1.Copy()
	c2.Namespace = "bar-cd"
	if c1 == c2 || c1.Namespace == c2.Namespace {
		t.Fatal("context not copied")
	}
}
