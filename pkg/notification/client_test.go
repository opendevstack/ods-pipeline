package notification

import (
	"bytes"
	"context"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const notificationJsonTemplate = `{
      "@type": "MessageCard",
      "@context": "http://schema.org/extensions",
      "themeColor": {{if eq .OverallStatus "Succeeded"}} "237b4b" {{else}} "c4314b" {{ end }},
      "summary": "{{.ODSContext.Project}} - ODS Pipeline Build finished with status {{.OverallStatus}}",
      "sections": [
        {
          "activityTitle": "ODS Pipeline Build finished with status {{.OverallStatus}}",
          "activitySubtitle": "On Project {{.ODSContext.Project}}",
          "activityImage": "https://avatars.githubusercontent.com/u/38974438?s=200&v=4",
          "facts": [
            {
              "name": "Component",
              "value": "{{.ODSContext.Component}}"
            },
            {
              "name": "Namespace",
              "value": "{{.ODSContext.Namespace}}"
            },
            {
              "name": "GitCommitSHA",
              "value": "{{.ODSContext.GitCommitSHA}}"
            },
            {
              "name": "GitRef",
              "value": "{{.ODSContext.GitRef}}"
            },
            {
              "name": "Version",
              "value": "{{.ODSContext.Version}}"
            },
            {
              "name": "Environment",
              "value": "{{.ODSContext.Environment}}"
            }
          ],
          "markdown": true
        }
      ],
      "potentialAction": [
        {
          "@type": "OpenUri",
          "name": "Go to PipelineRun",
          "targets": [
            {
              "os": "default",
              "uri": "{{.PipelineRunURL}}"
            }
          ]
        },
        {
          "@type": "OpenUri",
          "name": "Go to Git URL",
          "targets": [
            {
              "os": "default",
              "uri": "{{.ODSContext.GitURL}}"
            }
          ]
        }
        {{if .ODSContext.PullRequestBase}},
        {
          "@type": "OpenUri",
          "name": "Go to PR",
          "targets": [
            {
              "os": "default",
              "uri": "{{.ODSContext.PullRequestBase}}"
            }
          ]
        }
        {{end}}
      ]
    }`

func TestWebhookCall(t *testing.T) {
	runResult := PipelineRunResult{
		OverallStatus:  "Succeeded",
		PipelineRunURL: "https://localhost",
		ODSContext:     &pipelinectxt.ODSContext{},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedTemplate, err := template.New("expectedJson").Parse(notificationJsonTemplate)
		if err != nil {
			t.Fatalf("parsing jsonTemplate failed: %s", err)
		}
		want := bytes.NewBuffer([]byte{})
		err = expectedTemplate.Execute(want, runResult)
		if err != nil {
			t.Fatalf("error executing template: %s", err)
		}

		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("reading request body failed: %s", err)
		}

		if diff := cmp.Diff(want.String(), string(got)); diff != "" {
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
	err = webhookClient.CallWebhook(context.TODO(), runResult)
	if err != nil {
		t.Fatalf("call webhook failed: %s", err)
	}
}

func TestSkipNotification(t *testing.T) {
	runResult := PipelineRunResult{
		OverallStatus:  "None",
		PipelineRunURL: "https://localhost",
		ODSContext:     &pipelinectxt.ODSContext{},
	}
	allowedStatusValues := `["Failed","Succeeded"]`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("Server was called even though status '%s' not configured for notifications (%s)",
			runResult.OverallStatus, allowedStatusValues)
	}))
	defer ts.Close()

	kubernetesClient := kubernetes.TestClient{
		CMs: []*corev1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: NotificationConfigMap,
				},
				Data: map[string]string{
					UrlProperty:             ts.URL,
					MethodProperty:          "POST",
					ContentTypeProperty:     "application/json",
					NotifyOnStatusProperty:  allowedStatusValues,
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
	err = webhookClient.CallWebhook(context.TODO(), runResult)
	if err != nil {
		t.Fatalf("call webhook failed: %s", err)
	}
}
