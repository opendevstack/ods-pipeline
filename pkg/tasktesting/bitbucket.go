package tasktesting

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	kclient "k8s.io/client-go/kubernetes"
)

// BitbucketClientOrFatal returns a Bitbucket client, configured based on ConfigMap/Secret in the given namespace.
func BitbucketClientOrFatal(t *testing.T, c *kclient.Clientset, namespace string) *bitbucket.Client {
	bitbucketSecret, err := kubernetes.GetSecret(c, namespace, "ods-bitbucket-auth")
	if err != nil {
		t.Fatalf("could not get Bitbucket secret: %s", err)
	}
	bitbucketConfigMap, err := kubernetes.GetConfigMap(c, namespace, "ods-bitbucket")
	if err != nil {
		t.Fatalf("could not get Bitbucket config: %s", err)
	}

	bitbucketURL := strings.Replace(bitbucketConfigMap.Data["url"], "bitbucket-server-test.kind", "localhost", 1)
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: string(bitbucketSecret.Data["password"]),
		BaseURL:  bitbucketURL,
		Logger:   &logging.LeveledLogger{Level: logging.LevelDebug},
	})
	return bitbucketClient
}
