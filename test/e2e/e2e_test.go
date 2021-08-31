package e2e

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"knative.dev/pkg/apis"
)

func TestWebhookInterceptor(t *testing.T) {

	// Setup namespace to run tests in.
	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        tasktesting.StorageSourceDir,
			StorageCapacity:  tasktesting.StorageCapacity,
			StorageClassName: tasktesting.StorageClassName,
		},
	)

	// Cleanup namespace at the end.
	tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	defer tasktesting.TearDown(t, c, ns)

	// Create NodePort service which Bitbucket can post its webhook to.
	var nodePort int32 = 30950
	_, err := kubernetes.CreateNodePortService(
		c.KubernetesClientSet,
		"el-nodeport",
		map[string]string{"eventlistener": "ods-pipeline"},
		nodePort,
		8000,
		ns,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Initialize workspace with basic app.
	wsDir, err := tasktesting.InitWorkspace("source", "hello-world-app")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Workspace is in %s", wsDir)
	odsContext := tasktesting.SetupBitbucketRepo(
		t, c.KubernetesClientSet, ns, wsDir, tasktesting.BitbucketProjectKey,
	)

	// The webhook URL needs to be the address of the KinD control plane on the node port.
	ipAddress, err := kindControlPlaneIP()
	if err != nil {
		t.Fatal(err)
	}
	webhookURL := fmt.Sprintf("http://%s:%d", ipAddress, nodePort)
	t.Logf("Bitbucket webhook URL will be set to %s", webhookURL)

	// Create webhook in Bitbucket.
	webhookSecret, err := kubernetes.GetSecretKey(
		c.KubernetesClientSet, ns, "ods-bitbucket-webhook", "secret",
	)
	if err != nil {
		t.Fatalf("could not get Bitbucket webhook secret: %s", err)
	}
	bitbucketClient := tasktesting.BitbucketClientOrFatal(t, c.KubernetesClientSet, ns)
	_, err = bitbucketClient.WebhookCreate(
		odsContext.Project,
		odsContext.Repository,
		bitbucket.WebhookCreatePayload{
			Name:          "test",
			URL:           webhookURL,
			Active:        true,
			Events:        []string{"repo:refs_changed"},
			Configuration: bitbucket.WebhookConfiguration{Secret: webhookSecret},
		})
	if err != nil {
		t.Fatalf("could not create Bitbucket webhook: %s", err)
	}

	// Push a commit, which should trigger a webhook, which in turn should start a pipeline run.
	filename := "ods.yml"
	fileContent := `pipeline:
  tasks:
  - name: package-image
    taskRef:
      kind: ClusterTask
      name: ods-package-image
    workspaces:
    - name: source
      workspace: shared-workspace`

	err = ioutil.WriteFile(filepath.Join(wsDir, filename), []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("could not write file=%s: %s", filename, err)
	}
	requiredServices := []string{"ods-pipeline", "el-ods-pipeline", "el-nodeport"}
	serviceTimeout := time.Minute
	for _, serviceName := range requiredServices {
		t.Logf("Waiting %s for service %s to have ready pods ...\n", serviceTimeout, serviceName)
		err = waitForServiceToBeReady(t, c.KubernetesClientSet, ns, serviceName, serviceTimeout)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log("Sleeping for 10s to make it work - unsure why needed ...")
	time.Sleep(10 * time.Second)
	t.Log("Pushing file to Bitbucket ...")
	tasktesting.PushFileToBitbucketOrFatal(t, c.KubernetesClientSet, ns, wsDir, "master", "ods.yml")
	triggerTimeout := time.Minute
	t.Logf("Waiting %s for pipeline run to be triggered ...", triggerTimeout)
	pr, err := waitForPipelineRunToBeTriggered(c.TektonClientSet, ns, triggerTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Triggered pipeline run %s\n", pr.Name)
	runTimeout := 3 * time.Minute
	t.Logf("Waiting %s for pipeline run to succeed ...", runTimeout)
	gotReason, err := waitForPipelineRunToBeDone(c.TektonClientSet, ns, pr.Name, runTimeout)
	if err != nil {
		t.Fatal(err)
	}
	if gotReason != "Succeeded" {
		t.Logf("Want pipeline run reason to be 'Succeeded' but got '%s'", gotReason)
		logs, err := pipelineRunLogs(ns, pr.Name)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(logs)
		t.Fatal()
	}
}

func waitForServiceToBeReady(t *testing.T, clientset *k8s.Clientset, ns, name string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	svc, err := clientset.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for {
		time.Sleep(2 * time.Second)
		ready, reason, err := kubernetes.ServiceHasReadyPods(clientset, svc)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if ready {
			break
		}
		t.Logf("still waiting: %s ...", reason)
	}
	return nil
}

func waitForPipelineRunToBeTriggered(clientset *tekton.Clientset, ns string, timeout time.Duration) (*tektonv1beta1.PipelineRun, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var pipelineRunList *v1beta1.PipelineRunList
	for {
		time.Sleep(2 * time.Second)
		prs, err := clientset.TektonV1beta1().PipelineRuns(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if len(prs.Items) > 0 {
			pipelineRunList = prs
			break
		}
	}
	if len(pipelineRunList.Items) != 1 {
		return nil, errors.New("did not get exactly one pipeline run")
	}
	return &pipelineRunList.Items[0], nil
}

func waitForPipelineRunToBeDone(clientset *tekton.Clientset, ns, pr string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var reason string
	for {
		time.Sleep(2 * time.Second)
		pr, err := clientset.TektonV1beta1().PipelineRuns(ns).Get(ctx, pr, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		if pr.IsDone() {
			reason = pr.Status.GetCondition(apis.ConditionSucceeded).GetReason()
			break
		}
	}
	return reason, nil
}

func kindControlPlaneIP() (string, error) {
	stdout, stderr, err := command.Run(
		"docker",
		[]string{"inspect", "-f", "{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}", "kind-control-plane"},
	)
	if err != nil {
		return "", fmt.Errorf("could not get IP address of KinD control plane: %s, err: %s", string(stderr), err)
	}
	return strings.TrimSpace(string(stdout)), nil
}

func pipelineRunLogs(namespace, name string) (string, error) {
	if !tknInstalled() {
		return "", errors.New("tkn is not installed, cannot show logs")
	}
	stdout, stderr, err := command.Run(
		"tkn",
		[]string{"pr", "logs", "-n", namespace, name},
	)
	if err != nil {
		return "", fmt.Errorf("could not get logs of pipelinerun: %s, err: %s", string(stderr), err)
	}
	return string(stdout), nil
}

func tknInstalled() bool {
	_, err := exec.LookPath("tkn")
	return err == nil
}
