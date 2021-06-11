package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/internal/serverstub"
)

func TestProjectCreate(t *testing.T) {
	bitbucketClient := testClient(t, map[string]*serverstub.FakeResponse{
		"/rest/api/1.0/projects": {
			StatusCode: 201, Fixture: "bitbucket/project-create.json",
		},
	})
	p, err := bitbucketClient.ProjectCreate(ProjectCreatePayload{
		Key:  "PRJ",
		Name: "My Cool Project",
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Key != "PRJ" {
		t.Fatalf("got %s, want %s", p.Key, "PRJ")
	}
}
