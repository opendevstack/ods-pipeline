package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestTagCreate(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
		201, "bitbucket/tag-create.json",
	)

	tag, err := bitbucketClient.TagCreate("PRJ", "my-repo", TagCreatePayload{
		Name: "release-2.0.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if tag.ID != "release-2.0.0" {
		t.Fatalf("got %s, want %s", tag.ID, "release-2.0.0")
	}
}
