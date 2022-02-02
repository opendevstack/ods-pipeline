package notification

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/testfile"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookCall(t *testing.T) {
	notificationJsonTemplate := string(testfile.ReadFixture(t, "notification/teams-notification.tpl"))
	runResult := PipelineRunResult{
		PipelineRunName: "pipelinerun-release-v0-message-deadbeef",
		OverallStatus:   "Succeeded",
		PipelineRunURL:  "https://localhost",
		ODSContext: &pipelinectxt.ODSContext{
			Project:     "Project",
			GitRef:      "main",
			Environment: "dev",
			GitURL:      "https://localhost/vcs",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := string(testfile.ReadGolden(t, "notification/teams-notification.json"))

		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("reading request body failed: %s", err)
		}

		if diff := cmp.Diff(want, string(got)); diff != "" {
			t.Fatalf("Webhook expectation mismatch (-want +got):\n%s", diff)
		}
	}))
	defer ts.Close()
	kubernetesClient := kubernetes.TestClient{
		CMs: []*corev1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: NotificationConfigMap,
				},
				Data: map[string]string{
					EnabledProperty:         "true",
					UrlProperty:             ts.URL,
					MethodProperty:          "POST",
					ContentTypeProperty:     "application/json",
					NotifyOnStatusProperty:  `["Failed","Succeeded"]`,
					RequestTemplateProperty: notificationJsonTemplate,
				},
			},
		},
	}
	webhookClient, err := NewClient(ClientConfig{
		Namespace: "test",
	}, &kubernetesClient)
	if err != nil {
		t.Fatalf("constructing webhook client failed: %s", err)
	}
	ctxt := context.TODO()
	notificationConfig, err := webhookClient.ReadNotificationConfig(ctxt)
	if err != nil {
		t.Fatalf("Could not read notification config: %s", err)
	}

	err = webhookClient.CallWebhook(ctxt, notificationConfig, runResult)
	if err != nil {
		t.Fatalf("call webhook failed: %s", err)
	}
}

func TestSkipNotificationOnStatusMismatch(t *testing.T) {
	allowedStatusValues := []string{"Failed", "Succeeded"}
	webhookClient, err := NewClient(ClientConfig{
		Namespace: "test",
	}, nil)
	if err != nil {
		t.Fatalf("constructing webhook client failed: %s", err)
	}

	notificationConfig := &NotificationConfig{
		notifyOnStatus: allowedStatusValues,
	}

	status := "None"
	shouldNotify := webhookClient.ShouldNotify(notificationConfig, status)
	if shouldNotify {
		t.Fatalf("ShouldNotify was supposed to return false (status: %s, allowed: %s)",
			status, allowedStatusValues)
	}
}

func TestSkipNotificationOnNotificationsDisabled(t *testing.T) {
	webhookClient, err := NewClient(ClientConfig{
		Namespace: "test",
	}, nil)
	if err != nil {
		t.Fatalf("constructing webhook client failed: %s", err)
	}

	notificationConfig := &NotificationConfig{
		enabled: false,
	}

	if webhookClient.ShouldNotify(notificationConfig, "n/a") {
		t.Fatalf("ShouldNotify was supposed to return false when notifications are disabled")
	}
}
