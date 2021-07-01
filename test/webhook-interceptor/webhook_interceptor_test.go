package webhook_interceptor

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWebhookInterceptor(t *testing.T) {

	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard", // if using KinD, set it to "standard"
		},
	)

	_, err := kubernetes.CreateNodePortService(
		c.KubernetesClientSet,
		"el-test",
		map[string]string{"eventlistener": "ods-pipeline"},
		30950,
		8000,
		ns,
	)
	if err != nil {
		t.Fatal(err)
	}

	// tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	// defer tasktesting.TearDown(t, c, ns)

	wsDir, err := tasktesting.InitWorkspace("source", "hello-world-app")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Workspace is in %s", wsDir)

	bitbucketProjectKey := "ODSPIPELINETEST"
	odsContext := tasktesting.SetupBitbucketRepo(t, c.KubernetesClientSet, ns, wsDir, bitbucketProjectKey)

	// get webhook url
	// docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' kind-control-plane
	// might need https://kind.sigs.k8s.io/docs/user/configuration/#nodeport-with-port-mappings
	webhookURL := "http://172.18.0.3:30950"

	// create webhook setting
	bitbucketClient := tasktesting.BitbucketTestClient(t, c.KubernetesClientSet, ns)
	_, err = bitbucketClient.WebhookCreate(
		odsContext.Project,
		odsContext.Repository,
		bitbucket.WebhookCreatePayload{
			Name:          "test",
			URL:           webhookURL,
			Active:        true,
			Events:        []string{"repo:refs_changed"},
			Configuration: bitbucket.WebhookConfiguration{Secret: "test"}, // secret for Bitbucket
		})
	if err != nil {
		t.Fatalf("could not create Bitbucket webhook: %s", err)
	}
	// push a commit
	filename := "ods.yml"
	fileContent := `phases:
  build:
  - name: backend-build-go
    taskRef:
      kind: ClusterTask
      name: ods-build-go-v0-1-0
    workspaces:
    - name: source
      workspace: shared-workspace`
	// 	_, err = bitbucketClient.BrowseUpdate(
	// 		odsContext.Project,
	// 		odsContext.Repository,
	// 		filename,
	// 		bitbucket.BrowseUpdateParams{
	// 			Branch:         "master",
	// 			Message:        "initial commit",
	// 			SourceCommitId: "",
	// 			Content: strings.NewReader(fileContent),
	// 		},
	// 	)
	// 	if err != nil {
	// 		t.Fatalf("could not upload file to Bitbucket: %s", err)
	// 	}

	err = ioutil.WriteFile(filepath.Join(wsDir, filename), []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("could not write file=%s: %s", filename, err)
	}
	tasktesting.PushFileToBitbucketOrFatal(t, c.KubernetesClientSet, ns, wsDir, "master", "ods.yml")

	// figure out what the pipeline run is and wait for it to finish
	prs, err := c.TektonClientSet.TektonV1beta1().PipelineRuns(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("could not get pipeline runs: %s", err)
	}
	for _, pr := range prs.Items {
		fmt.Println(pr.Name)
	}

	// check it is a success
}
