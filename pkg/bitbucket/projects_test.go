package bitbucket

import (
	"testing"

	"github.com/opendevstack/ods-pipeline/test/testserver"
)

func TestProjectCreate(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(t, srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects",
		201, "bitbucket/project-create.json",
	)

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
