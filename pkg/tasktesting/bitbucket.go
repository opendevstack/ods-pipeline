package tasktesting

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/installation"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	kclient "k8s.io/client-go/kubernetes"
)

const (
	BitbucketProjectKey = "ODSPIPELINETEST"
	BitbucketKinDHost   = "ods-test-bitbucket-server.kind"
)

// BitbucketClientOrFatal returns a Bitbucket client, configured based on ConfigMap/Secret in the given namespace.
func BitbucketClientOrFatal(t *testing.T, c *kclient.Clientset, namespace string) *bitbucket.Client {
	bcc, err := installation.NewBitbucketClientConfig(
		c, namespace, &logging.LeveledLogger{Level: logging.LevelDebug},
	)
	if err != nil {
		t.Fatalf("could not create Bitbucket client config: %s", err)
	}
	bcc.BaseURL = strings.Replace(bcc.BaseURL, BitbucketKinDHost, "localhost", 1)
	bc, err := bitbucket.NewClient(bcc)
	if err != nil {
		t.Fatalf("bitbucket client: %s", err)
	}
	return bc
}
