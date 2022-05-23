package bitbucket

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestBuildStatusCreate(t *testing.T) {
	sha := "56625c80087b034847001d22502063adae9759f2"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	payload := BuildStatusCreatePayload{
		State:       "INPROGRESS",
		Key:         sha,
		Name:        sha,
		URL:         "http://acme.org/pipeline-run",
		Description: "ODS Pipeline Build",
	}

	srv.EnqueueResponse(
		t, "/rest/build-status/1.0/commits/"+sha,
		204, "",
	)
	err := bitbucketClient.BuildStatusCreate(sha, payload)
	if err != nil {
		t.Fatal(err)
	}

	srv.EnqueueResponse(
		t, "/rest/build-status/1.0/commits/"+sha,
		500, "bitbucket/error.txt",
	)
	err = bitbucketClient.BuildStatusCreate(sha, payload)
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}

func TestBuildStatusList(t *testing.T) {
	sha := "56625c80087b034847001d22502063adae9759f2"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/build-status/1.0/commits/"+sha,
		200, "bitbucket/build-status-list.json",
	)
	l, err := bitbucketClient.BuildStatusList(sha)
	if err != nil {
		t.Fatal(err)
	}
	if l.Size != 1 {
		t.Fatalf("got %d, want %d", l.Size, 1)
	}

	srv.EnqueueResponse(
		t, "/rest/build-status/1.0/commits/"+sha,
		500, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.BuildStatusList(sha)
	want := "status code: 500, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}

	srv.EnqueueResponse(
		t, "/rest/build-status/1.0/commits/"+sha,
		200, "bitbucket/error.txt",
	)
	_, err = bitbucketClient.BuildStatusList(sha)
	want = "unmarshal: invalid character 'e' looking for beginning of value. status code: 200, body: error description"
	if strings.TrimSpace(err.Error()) != want {
		t.Fatalf("got error: %s. want: %s", err, want)
	}
}
