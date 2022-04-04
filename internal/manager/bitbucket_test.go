package manager

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
)

func TestGetRepoNames(t *testing.T) {
	c := &bitbucket.TestClient{
		Repos: []bitbucket.Repo{{Name: "a"}, {Name: "b"}},
	}
	got, err := GetRepoNames(c, "PROJ")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a", "b"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("expected (-want +got):\n%s", diff)
	}
}

func TestGetCommitSHA(t *testing.T) {
	c := &bitbucket.TestClient{
		Commits: []bitbucket.Commit{{ID: "a"}, {ID: "b"}},
	}
	got, err := getCommitSHA(c, "PROJ", "repo", "branch")
	if err != nil {
		t.Fatal(err)
	}
	want := "a"
	if want != got {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}
