package tasktesting

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/installation"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	kclient "k8s.io/client-go/kubernetes"
)

const (
	NexusKinDHost = "ods-test-nexus.kind"
)

// NexusClientOrFatal returns a Nexus client, configured based on ConfigMap/Secret in the given namespace.
func NexusClientOrFatal(t *testing.T, c *kclient.Clientset, namespace string) *nexus.Client {
	ncc, err := installation.NewNexusClientConfig(
		c, namespace, &logging.LeveledLogger{Level: logging.LevelDebug},
	)
	if err != nil {
		t.Fatalf("could not create Nexus client config: %s", err)
	}
	ncc.BaseURL = strings.Replace(ncc.BaseURL, NexusKinDHost, "localhost", 1)

	nexusClient, err := nexus.NewClient(ncc)
	if err != nil {
		t.Fatal(err)
	}
	return nexusClient
}
