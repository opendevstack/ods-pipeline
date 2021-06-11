package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/internal/serverstub"
)

func TestTagCreate(t *testing.T) {
	bitbucketClient := testClient(t, map[string]*serverstub.FakeResponse{
		"/rest/api/1.0/projects/PRJ/repos/my-repo/tags": {
			StatusCode: 201, Fixture: "bitbucket/tag-create.json",
		},
	})
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
