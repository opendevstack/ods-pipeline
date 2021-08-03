package pipelinectxt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAssemble(t *testing.T) {
	c := &ODSContext{
		Namespace:    "foo-cd",
		GitFullRef:   "refs/heads/master",
		GitURL:       "https://github.com/opendevstack/ods-pipeline.git",
		GitCommitSHA: "abcdef",
	}
	err := c.Assemble(".")
	if err != nil {
		t.Fatal(err)
	}
	wantContext := &ODSContext{
		Namespace:       "foo-cd",
		Project:         "foo",
		Repository:      "pipelinectxt",
		Component:       "pipelinectxt",
		GitCommitSHA:    "abcdef",
		GitFullRef:      "refs/heads/master",
		GitRef:          "master",
		GitURL:          "https://github.com/opendevstack/ods-pipeline.git",
		Version:         "WIP",
		Environment:     "",
		PullRequestBase: "",
		PullRequestKey:  "",
	}
	if diff := cmp.Diff(wantContext, c); diff != "" {
		t.Fatalf("context mismatch (-want +got):\n%s", diff)
	}
}
