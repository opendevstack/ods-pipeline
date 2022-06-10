package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestCommitGet(t *testing.T) {
	sha := "abcdef0123abcdef4567abcdef8987abcdef6543"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(t, srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits/"+sha,
		200, "bitbucket/commit-get.json",
	)

	c, err := bitbucketClient.CommitGet("myproject", "my-repo", sha)
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != sha {
		t.Fatalf("got %s, want %s", c.ID, sha)
	}
}

func TestCommitList(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(t, srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits",
		200, "bitbucket/commit-list.json",
	)

	l, err := bitbucketClient.CommitList("myproject", "my-repo", CommitListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if l.Size != 1 {
		t.Fatalf("got %d, want %d", l.Size, 1)
	}
}

func TestCommitPullRequestList(t *testing.T) {
	sha := "abcdef0123abcdef4567abcdef8987abcdef6543"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(t, srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits/"+sha+"/pull-requests",
		200, "bitbucket/commit-pull-request-list.json",
	)

	l, err := bitbucketClient.CommitPullRequestList("myproject", "my-repo", sha)
	if err != nil {
		t.Fatal(err)
	}
	if l.Size != 1 {
		t.Fatalf("got %d, want %d", l.Size, 1)
	}
}
