package bitbucket

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/serverstub"
)

func TestRawGet(t *testing.T) {
	at := "refs/heads/master"
	bitbucketClient := testClient(t, map[string]*serverstub.FakeResponse{
		"/projects/PRJ/repos/my-repo/raw/example.txt": {
			StatusCode: 200, Fixture: "bitbucket/example.txt",
		},
	})
	r, err := bitbucketClient.RawGet(
		"PRJ", "my-repo", "example.txt", at,
	)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(string(r))
	if got != "hello world" {
		t.Fatalf("got %s, want %s", got, "hello world")
	}
}
