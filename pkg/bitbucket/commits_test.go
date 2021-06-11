package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/internal/serverstub"
)

func TestCommitGet(t *testing.T) {
	sha := "abcdef0123abcdef4567abcdef8987abcdef6543"
	bitbucketClient := testClient(t, map[string]*serverstub.FakeResponse{
		"/rest/api/1.0/projects/myproject/repos/myrepo/commits/" + sha: {
			StatusCode: 200, Fixture: "bitbucket/commit-get.json",
		},
	})
	c, err := bitbucketClient.CommitGet("myproject", "myrepo", sha)
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != sha {
		t.Fatalf("got %s, want %s", c.ID, sha)
	}
}

func TestCommitList(t *testing.T) {
	bitbucketClient := testClient(t, map[string]*serverstub.FakeResponse{
		"/rest/api/1.0/projects/myproject/repos/myrepo/commits": {
			StatusCode: 200, Fixture: "bitbucket/commit-list.json",
		},
	})
	l, err := bitbucketClient.CommitList("myproject", "myrepo", CommitListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if l.Size != 1 {
		t.Fatalf("got %d, want %d", l.Size, 1)
	}
}
