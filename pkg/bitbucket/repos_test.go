package bitbucket

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestRepoCreate(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	payload := RepoCreatePayload{
		Name: "my-repo",
	}

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos",
		201, "bitbucket/repo-create.json",
	)
	r, err := bitbucketClient.RepoCreate("PRJ", payload)
	if err != nil {
		t.Fatal(err)
	}
	if r.Slug != "my-repo" {
		t.Fatalf("got %s, want %s", r.Slug, "my-repo")
	}

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos",
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.RepoCreate("PRJ", payload)
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos",
		201, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.RepoCreate("PRJ", payload)
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 201, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}

func TestRepoList(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos",
		200, "bitbucket/repo-list.json",
	)
	l, err := bitbucketClient.RepoList("PRJ")
	if err != nil {
		t.Fatal(err)
	}
	if l.Size != 1 {
		t.Fatalf("got %d, want %d", l.Size, 1)
	}

	// check status code error is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos",
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.RepoList("PRJ")
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	// check incorrect body is handled properly
	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos",
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.RepoList("PRJ")
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}
