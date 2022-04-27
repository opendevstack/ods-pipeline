package bitbucket

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestCommitGet(t *testing.T) {
	sha := "abcdef0123abcdef4567abcdef8987abcdef6543"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

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

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits/"+sha,
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.CommitGet("myproject", "my-repo", sha)
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits/"+sha,
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.CommitGet("myproject", "my-repo", sha)
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}

func TestCommitList(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

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

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits",
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.CommitList("myproject", "my-repo", CommitListParams{})
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits",
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.CommitList("myproject", "my-repo", CommitListParams{})
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}

func TestCommitPullRequestList(t *testing.T) {
	sha := "abcdef0123abcdef4567abcdef8987abcdef6543"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

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

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits/"+sha+"/pull-requests",
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.CommitPullRequestList("myproject", "my-repo", sha)
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/myproject/repos/my-repo/commits/"+sha+"/pull-requests",
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.CommitPullRequestList("myproject", "my-repo", sha)
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}
