package tasktesting

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	kclient "k8s.io/client-go/kubernetes"
)

const (
	BitbucketProjectKey = "ODSPIPELINETEST"
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

	bitbucketURL := strings.Replace(bitbucketConfigMap.Data["url"], "ods-test-bitbucket-server.kind", "localhost", 1)
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: string(bitbucketSecret.Data["password"]),
		BaseURL:  bitbucketURL,
		Logger:   &logging.LeveledLogger{Level: logging.LevelDebug},
	})
	return bitbucketClient
}

func CheckBitbucketBuildStatus(t *testing.T, c *bitbucket.Client, gitCommit, wantBuildStatus string) {
	buildStatusPage, err := c.BuildStatusList(gitCommit)
	if err != nil {
		t.Fatal(err)
	}
	if buildStatusPage == nil || len(buildStatusPage.Values) == 0 {
		t.Fatal("no build status found")
	}
	buildStatus := buildStatusPage.Values[0]
	if buildStatus.State != wantBuildStatus {
		t.Fatalf("Got: %s, want: %s", buildStatus.State, wantBuildStatus)
	}

}
