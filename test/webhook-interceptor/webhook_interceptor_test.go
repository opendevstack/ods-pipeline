package webhook_interceptor

import (
	"testing"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestWebhookInterceptor(t *testing.T) {

	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard", // if using KinD, set it to "standard"
		},
	)

	tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	defer tasktesting.TearDown(t, c, ns)

	wsDir, err := tasktesting.InitWorkspace("source", "hello-world-app")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Workspace is in %s", wsDir)

	bitbucketProjectKey := "ODSPIPELINETEST"
	odsContext := tasktesting.SetupBitbucketRepo(t, c.KubernetesClientSet, ns, wsDir, bitbucketProjectKey)

	// create webhook setting
	bitbucketClient := tasktesting.BitbucketTestClient(t, c, ns)
	_, err = bitbucketClient.WebhookCreate(
		odsContext.Project,
		odsContext.Repository,
		bitbucket.WebhookCreatePayload{
			Name:          "test",
			URL:           "", // URL of event listener
			Active:        true,
			Events:        []string{"repo:refs_changed"},
			Configuration: bitbucket.WebhookConfiguration{Secret: ""}, // secret for Bitbucket
		})
	if err != nil {
		t.Fatalf("could not create Bitbucket webhook: %s", err)
	}
	// push a commit
	// see https://docs.atlassian.com/bitbucket-server/rest/6.4.0/bitbucket-rest.html#idp188

	// figure out what the pipeline run is and wait for it to finish

	// check it is a success
}
