package notification

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/testfile"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

func TestWebhookCall(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := string(testfile.ReadGolden(t, "notification/teams-notification.json"))

		got, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("reading request body failed: %s", err)
		}

		if diff := cmp.Diff(want, string(got)); diff != "" {
			t.Fatalf("Webhook expectation mismatch (-want +got):\n%s", diff)
		}
	}))
	defer ts.Close()

	notificationJsonTemplate := string(testfile.ReadFixture(t, "notification/teams-notification.tpl"))
	payloadTemplate, err := template.New("notificiationPayload").Parse(notificationJsonTemplate)
	if err != nil {
		t.Fatalf("parsing json Template fixture failed: %v", err)
	}
	notificationConfig := Config{
		Enabled:        true,
		URL:            ts.URL,
		Method:         "POST",
		ContentType:    "application/json",
		NotifyOnStatus: []string{"Failed", "Succeeded"},
		Template:       payloadTemplate,
	}
	runResult := PipelineRunResult{
		PipelineRunName: "pipelinerun-release-v0-message-deadbeef",
		OverallStatus:   "Succeeded",
		PipelineRunURL:  "https://localhost",
		ODSContext: &pipelinectxt.ODSContext{
			Project: "Project",
			GitRef:  "main",
			GitURL:  "https://localhost/vcs",
		},
	}

	webhookClient, err := NewClient(ClientConfig{
		Namespace:          "test",
		NotificationConfig: &notificationConfig,
	})
	if err != nil {
		t.Fatalf("constructing webhook client failed: %s", err)
	}

	ctxt := context.TODO()
	err = webhookClient.CallWebhook(ctxt, runResult)
	if err != nil {
		t.Fatalf("call webhook failed: %s", err)
	}
}

func TestSkipNotificationOnStatusMismatch(t *testing.T) {
	allowedStatusValues := []string{"Failed", "Succeeded"}
	notificationConfig := &Config{
		NotifyOnStatus: allowedStatusValues,
	}
	webhookClient, err := NewClient(ClientConfig{
		Namespace:          "test",
		NotificationConfig: notificationConfig,
	})
	if err != nil {
		t.Fatalf("constructing webhook client failed: %s", err)
	}

	status := "None"
	shouldNotify := webhookClient.ShouldNotify(status)
	if shouldNotify {
		t.Fatalf("ShouldNotify was supposed to return false (status: %s, allowed: %s)",
			status, allowedStatusValues)
	}
}

func TestSkipNotificationOnNotificationsDisabled(t *testing.T) {
	notificationConfig := &Config{
		Enabled: false,
	}
	webhookClient, err := NewClient(ClientConfig{
		Namespace:          "test",
		NotificationConfig: notificationConfig,
	})
	if err != nil {
		t.Fatalf("constructing webhook client failed: %s", err)
	}

	if webhookClient.ShouldNotify("n/a") {
		t.Fatalf("ShouldNotify was supposed to return false when notifications are disabled")
	}
}
