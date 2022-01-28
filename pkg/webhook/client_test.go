package webhook

import (
	"bytes"
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
)

const jsonTemplate = `{
    "@context": "https://schema.org/extensions",
    "@type": "MessageCard",
    "themeColor": "c60000",
    "title": "ODS Pipeline Build finished",
    "text": "ODS Pipeline run finished with status {{.OverallStatus}}!",
    "potentialAction": [
        {
            "@type": "OpenUri",
            "name": "Learn More",
            "targets": [
                { "os": "default", "uri":{{.PipelineRunURL}}" }
            ]
        }
    ]
}`

func TestWebhookCall(t *testing.T) {
	runResult := PipelineRunResult{
		OverallStatus:  "SUCCESS",
		PipelineRunURL: "https://localhost",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedTemplate, err := template.New("expectedJson").Parse(jsonTemplate)
		if err != nil {
			t.Fatalf("parsing jsonTemplate failed: %s", err)
		}
		want := bytes.NewBuffer([]byte{})
		err = expectedTemplate.Execute(want, runResult)

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
					Name: webhookConfigMap,
				},
				Data: map[string]string{
					urlProperty:             ts.URL,
					methodProperty:          "POST",
					contentTypeProperty:     "application/json",
					requestTemplateProperty: jsonTemplate,
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
