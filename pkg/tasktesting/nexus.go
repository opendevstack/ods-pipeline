package tasktesting

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/ods-pipeline/internal/installation"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
	"github.com/opendevstack/ods-pipeline/pkg/nexus"
	kclient "k8s.io/client-go/kubernetes"
)

const (
	NexusKinDHost    = "ods-test-nexus.kind"
	NexusKinDTLSHost = "ods-test-nexus-tls.kind"
)

// NexusClientOrFatal returns a Nexus client, configured based on ConfigMap/Secret in the given namespace.
func NexusClientOrFatal(t *testing.T, c *kclient.Clientset, namespace string, privateCert bool) *nexus.Client {
	ncc, err := installation.NewNexusClientConfig(
		c, namespace, &logging.LeveledLogger{Level: logging.LevelDebug},
	)
	if err != nil {
		t.Fatalf("could not create Nexus client config: %s", err)
	}
	ncc.BaseURL = strings.Replace(ncc.BaseURL, NexusKinDHost, "localhost", 1)
	ncc.BaseURL = strings.Replace(ncc.BaseURL, NexusKinDTLSHost, "localhost", 1)
	if true {
		ncc.CertFile = filepath.Join(projectpath.Root, PrivateCertFile)
	}
	nexusClient, err := nexus.NewClient(ncc)
	if err != nil {
		t.Fatal(err)
	}
	return nexusClient
}
