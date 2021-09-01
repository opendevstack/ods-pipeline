package tasktesting

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	kclient "k8s.io/client-go/kubernetes"
)

// NexusClientOrFatal returns a Nexus client, configured based on ConfigMap/Secret in the given namespace.
func NexusClientOrFatal(t *testing.T, c *kclient.Clientset, namespace string) *nexus.Client {
	nexusSecret, err := kubernetes.GetSecret(c, namespace, "ods-nexus-auth")
	if err != nil {
		t.Fatalf("could not get Nexus secret: %s", err)
	}
	nexusConfigMap, err := kubernetes.GetConfigMap(c, namespace, "ods-nexus")
	if err != nil {
		t.Fatalf("could not get Nexus config: %s", err)
	}

	nexusURL := strings.Replace(nexusConfigMap.Data["url"], "ods-test-nexus.kind", "localhost", 1)
	nexusClient, err := nexus.NewClient(&nexus.ClientConfig{
		BaseURL:  nexusURL,
		Username: string(nexusSecret.Data["username"]),
		Password: string(nexusSecret.Data["password"]),
		Logger:   &logging.LeveledLogger{Level: logging.LevelDebug},
	})
	if err != nil {
		t.Fatal(err)
	}
	return nexusClient
}
