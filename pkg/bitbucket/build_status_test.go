package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/internal/serverstub"
)

func TestBuildStatusCreate(t *testing.T) {
	sha := "56625c80087b034847001d22502063adae9759f2"
	bitbucketClient := testClient(t, map[string]*serverstub.FakeResponse{
		"/rest/build-status/1.0/commits/" + sha: {
			StatusCode: 204, Fixture: "",
		},
	})
	err := bitbucketClient.BuildStatusCreate(sha, BuildStatusPostPayload{
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
