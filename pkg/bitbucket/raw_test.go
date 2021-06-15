package bitbucket

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestRawGet(t *testing.T) {
	at := "refs/heads/master"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/projects/PRJ/repos/my-repo/raw/example.txt",
		200, "bitbucket/example.txt",
	)

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
