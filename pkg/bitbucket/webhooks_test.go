package bitbucket

import (
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestWebhookCreate(t *testing.T) {
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(t, srv.Server.URL)

	srv.EnqueueResponse(
		t, "/rest/api/1.0/projects/PRJ/repos/my-repo/webhooks",
		201, "bitbucket/webhook-create.json",
	)

	w, err := bitbucketClient.WebhookCreate(
		"PRJ", "my-repo",
		WebhookCreatePayload{
			Name:   "Webhook Name",
			Events: []string{"repo:refs_changed", "repo:modified"},
			Configuration: WebhookConfiguration{
				Secret: "password",
			},
			URL:    "http://example.com",
			Active: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if w.ID != 10 {
		t.Fatalf("got %d, want %d", w.ID, 10)
	}
}
