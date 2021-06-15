package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestRepoCreate(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos",
		201, "bitbucket/repo-create.json",
	)

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
