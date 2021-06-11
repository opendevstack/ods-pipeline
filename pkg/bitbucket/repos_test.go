package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/internal/serverstub"
)

func TestRepoCreate(t *testing.T) {
	bitbucketClient := testClient(t, map[string]*serverstub.FakeResponse{
		"/rest/api/1.0/projects/PRJ/repos": {
			StatusCode: 201, Fixture: "bitbucket/repo-create.json",
		},
	})
	r, err := bitbucketClient.RepoCreate("PRJ", RepoCreatePayload{
		Name: "my-repo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Slug != "my-repo" {
		t.Fatalf("got %s, want %s", r.Slug, "my-repo")
	}
}
