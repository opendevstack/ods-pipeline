package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestBuildStatusCreate(t *testing.T) {
	sha := "56625c80087b034847001d22502063adae9759f2"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/build-status/1.0/commits/"+sha,
		204, "",
	)

	err := bitbucketClient.BuildStatusCreate(sha, BuildStatusCreatePayload{
		State:       "INPROGRESS",
		Key:         sha,
		Name:        sha,
		URL:         "http://acme.org/pipeline-run",
		Description: "ODS Pipeline Build",
	})
	if err != nil {
		t.Fatal(err)
	}
}
