package pipelinectxt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAssemble(t *testing.T) {
	c := &ODSContext{
		Namespace:    "foo-cd",
		GitFullRef:   "refs/heads/master",
		GitURL:       "https://example.bitbucket.com/scm/ODS/ods-pipeline.git",
		GitCommitSHA: "abcdef",
	}
	err := c.Assemble(".")
	if err != nil {
		t.Fatal(err)
	}
	wantContext := &ODSContext{
		Namespace:       "foo-cd",
		Project:         "ods",
		Repository:      "ods-pipeline",
		Component:       "pipeline",
		GitCommitSHA:    "abcdef",
		GitFullRef:      "refs/heads/master",
		GitRef:          "master",
		GitURL:          "https://example.bitbucket.com/scm/ODS/ods-pipeline.git",
		Version:         "WIP",
		Environment:     "",
		PullRequestBase: "",
		PullRequestKey:  "",
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
