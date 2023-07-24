package e2e

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/internal/kubernetes"
	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	"github.com/opendevstack/ods-pipeline/pkg/tasktesting"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"knative.dev/pkg/apis"
)

var outsideKindFlag = flag.Bool("outside-kind", false, "Whether to continue if not in KinD cluster")
var privateCertFlag = flag.Bool("private-cert", false, "Whether to run tests using a private cert")

func TestE2E(t *testing.T) {
	tasktesting.CheckCluster(t, *outsideKindFlag)
	tasktesting.CheckServices(t, []tasktesting.Service{
		tasktesting.Bitbucket, tasktesting.Nexus,
	})

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
		"ods-pm-nodeport", // NodePort for ODS Pipeline Manager
		map[string]string{
			"app.kubernetes.io/name":     "ods-pipeline",
			"app.kubernetes.io/instance": "ods-pipeline",
		},
		nodePort,
		8080,
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
		t, c.KubernetesClientSet, ns, wsDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
	)

	// The webhook URL needs to be the address of the KinD control plane on the node port.
	ipAddress, err := kindControlPlaneIP()
	if err != nil {
		t.Fatal(err)
	}
	webhookURL := fmt.Sprintf("http://%s:%d/bitbucket", ipAddress, nodePort)
	t.Logf("Bitbucket webhook URL will be set to %s", webhookURL)

	// Create webhook in Bitbucket.
	webhookSecret, err := kubernetes.GetSecretKey(
		c.KubernetesClientSet, ns, "ods-bitbucket-webhook", "secret",
	)
	if err != nil {
		t.Fatalf("could not get Bitbucket webhook secret: %s", err)
	}
	bitbucketClient := tasktesting.BitbucketClientOrFatal(t, c.KubernetesClientSet, ns, *privateCertFlag)
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
	filename := "ods.yaml"
	fileContent := `
pipelines:
  - tasks:
      - name: package-image
        taskRef:
          kind: Task
          name: ods-package-image
        workspaces:
          - name: source
            workspace: shared-workspace`

	err = os.WriteFile(filepath.Join(wsDir, filename), []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("could not write file=%s: %s", filename, err)
	}
	requiredService := "ods-pipeline"
	serviceTimeout := time.Minute
	t.Logf("Waiting %s for service %s to have ready pods ...\n", serviceTimeout, requiredService)
	err = waitForServiceToBeReady(t, c.KubernetesClientSet, ns, requiredService, serviceTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Pushing file to Bitbucket ...")
	tasktesting.PushFileToBitbucketOrFatal(t, c.KubernetesClientSet, ns, wsDir, "master:feature/test-branch", "ods.yaml")
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
	var svc *corev1.Service
	for {
		if svc == nil {
			s, err := clientset.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			svc = s
		}
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

	var pipelineRunList *tektonv1beta1.PipelineRunList
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
	stdout, stderr, err := command.RunBuffered(
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
	stdout, stderr, err := command.RunBuffered(
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
