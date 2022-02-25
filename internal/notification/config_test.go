package notification

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"text/template"
	"text/template/parse"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/testfile"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestReadNotificationConfigFromConfigMap(t *testing.T) {
	notificationJsonTemplate := string(testfile.ReadFixture(t, "notification/teams-notification.tpl"))
	payloadTemplate, err := template.New("requestTemplate").Parse(notificationJsonTemplate)
	if err != nil {
		t.Fatalf("failed to parse Template from fixture: %v", err)
	}

	want := &Config{
		Enabled:        true,
		URL:            "https://localhost",
		Method:         "POST",
		ContentType:    "application/json",
		NotifyOnStatus: []string{"Failed", "Succeeded"},
		Template:       payloadTemplate,
	}

	statusValues, err := json.Marshal(want.NotifyOnStatus)
	if err != nil {
		t.Fatalf("failed to marshal NotifyOnStatus to json: %v", err)
	}
	kubernetesClient := kubernetes.TestClient{
		CMs: []*corev1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: configMapName,
				},
				Data: map[string]string{
					enabledProperty:         strconv.FormatBool(want.Enabled),
					urlProperty:             want.URL,
					methodProperty:          want.Method,
					contentTypeProperty:     want.ContentType,
					notifyOnStatusProperty:  string(statusValues),
					requestTemplateProperty: notificationJsonTemplate,
				},
			},
		},
	}

	got, err := ReadConfigFromConfigMap(context.TODO(), &kubernetesClient)
	if err != nil {
		t.Fatalf("failed to read notification config: %v", err)
	}

	// need to ignore unexported fields template AST nodes to allow comparing templates
	ignoreUnexported := cmpopts.IgnoreUnexported(parse.Tree{}, template.Template{}, parse.ListNode{},
		parse.TextNode{}, parse.ActionNode{}, parse.PipeNode{}, parse.BranchNode{}, parse.BranchNode{},
		parse.CommandNode{}, parse.IdentifierNode{}, parse.FieldNode{}, parse.StringNode{})
	if diff := cmp.Diff(want, got, cmp.Options{ignoreUnexported}); diff != "" {
		t.Fatalf("Webhook expectation mismatch (-want +got):\n%s", diff)
	}
}
