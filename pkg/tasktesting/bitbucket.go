package tasktesting

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/installation"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	kclient "k8s.io/client-go/kubernetes"
)

const (
	BitbucketProjectKey  = "ODSPIPELINETEST"
	BitbucketKinDHost    = "ods-test-bitbucket-server.kind"
	BitbucketTLSKinDHost = "ods-test-bitbucket-server-tls.kind"
)

// BitbucketClientOrFatal returns a Bitbucket client, configured based on ConfigMap/Secret in the given namespace.
func BitbucketClientOrFatal(t *testing.T, c *kclient.Clientset, namespace string, privateCert bool) *bitbucket.Client {
	var privateCertPath string
	if privateCert {
		privateCertPath = filepath.Join(projectpath.Root, PrivateCertFile)
	}
	bcc, err := installation.NewBitbucketClientConfig(
		c, namespace, &logging.LeveledLogger{Level: logging.LevelDebug}, privateCertPath,
	)
	if err != nil {
		t.Fatalf("could not create Bitbucket client config: %s", err)
	}
	bcc.BaseURL = strings.Replace(bcc.BaseURL, BitbucketKinDHost, "localhost", 1)
	bcc.BaseURL = strings.Replace(bcc.BaseURL, BitbucketTLSKinDHost, "localhost", 1)
	bc, err := bitbucket.NewClient(bcc)
	if err != nil {
		t.Fatalf("bitbucket client: %s", err)
	}
	return bc
}
