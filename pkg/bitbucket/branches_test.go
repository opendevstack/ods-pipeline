package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestBranchList(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/branches",
		200, "bitbucket/branch-list.json",
	)

	l, err := bitbucketClient.BranchList("myproject", "my-repo", BranchListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if l.Size != 1 {
		t.Fatalf("got %d, want %d", l.Size, 1)
	}
}
