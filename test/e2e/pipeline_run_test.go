package e2e

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/internal/kubernetes"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	"github.com/opendevstack/ods-pipeline/internal/tasktesting"
	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	"github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8s "k8s.io/client-go/kubernetes"
	"knative.dev/pkg/apis"
)

func TestPipelineRun(t *testing.T) {
	k8sClient := newK8sClient(t)
	// Create NodePort service which Bitbucket can post its webhook to.
	var nodePort int32 = 30950
	_, err := createNodePortService(
		k8sClient,
		"ods-pm-nodeport", // NodePort for ODS Pipeline Manager
		map[string]string{
			"app.kubernetes.io/name":     "ods-pipeline",
			"app.kubernetes.io/instance": "ods-pipeline",
		},
		nodePort,
		8080,
		namespaceConfig.Name,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Initialize workspace with basic app.
	workspaceSourceDirectory := filepath.Join(
		projectpath.Root, "test", testdataWorkspacesPath, "hello-world-app",
	)
	wsDir, wsDirCleanupFunc, err := tektontaskrun.SetupWorkspaceDir(workspaceSourceDirectory)
	defer wsDirCleanupFunc()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Workspace is in %s", wsDir)
	odsContext := tasktesting.SetupBitbucketRepo(
		t, k8sClient, namespaceConfig.Name, wsDir, tasktesting.BitbucketProjectKey, *privateCertFlag,
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
		k8sClient, namespaceConfig.Name, "ods-bitbucket-webhook", "secret",
	)
	if err != nil {
		t.Fatalf("could not get Bitbucket webhook secret: %s", err)
	}
	bitbucketClient := tasktesting.BitbucketClientOrFatal(t, k8sClient, namespaceConfig.Name, *privateCertFlag)
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
  - name: hello-world
    taskSpec:
      steps:
      - name: message
        image: busybox
        script: |
          echo "hello world"
        workingDir: $(workspaces.source.path)
      workspaces:
      - name: source
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
	err = waitForServiceToBeReady(t, k8sClient, namespaceConfig.Name, requiredService, serviceTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Pushing file to Bitbucket ...")
	tasktesting.PushFileToBitbucketOrFatal(t, k8sClient, namespaceConfig.Name, wsDir, "master:feature/test-branch", "ods.yaml")
	triggerTimeout := time.Minute
	tektonClient := newTektonClient(t)
	t.Logf("Waiting %s for pipeline run to be triggered ...", triggerTimeout)
	pr, err := waitForPipelineRunToBeTriggered(tektonClient, namespaceConfig.Name, triggerTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Triggered pipeline run %s\n", pr.Name)
	runTimeout := 3 * time.Minute
	t.Logf("Waiting %s for pipeline run to succeed ...", runTimeout)
	gotReason, err := waitForPipelineRunToBeDone(tektonClient, namespaceConfig.Name, pr.Name, runTimeout)
	if err != nil {
		t.Fatal(err)
	}
	if gotReason != "Succeeded" {
		t.Logf("Want pipeline run reason to be 'Succeeded' but got '%s'", gotReason)
		logs, err := pipelineRunLogs(namespaceConfig.Name, pr.Name)
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
		ready, reason, err := serviceHasReadyPods(clientset, svc)
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
		[]string{"inspect", "-f", "{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}", tektontaskrun.KinDName + "-control-plane"},
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

func createNodePortService(clientset k8s.Interface, name string, selectors map[string]string, port, targetPort int32, namespace string) (*corev1.Service, error) {
	log.Printf("Create node port service %s", name)
	svc, err := clientset.CoreV1().Services(namespace).Create(context.TODO(),
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Spec: corev1.ServiceSpec{
				ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
				Ports: []corev1.ServicePort{
					{
						Name:       fmt.Sprintf("%d-%d", port, targetPort),
						NodePort:   port,
						Port:       port,
						Protocol:   corev1.ProtocolTCP,
						TargetPort: intstr.FromInt(int(targetPort)),
					},
				},
				Selector:        selectors,
				SessionAffinity: corev1.ServiceAffinityNone,
				Type:            corev1.ServiceTypeNodePort,
			},
		}, metav1.CreateOptions{})

	return svc, err
}

// serviceHasReadyPods returns false if no pod is assigned to given service
// or if one or more pods are not "Running"
// or one or more of any pods containers are not "ready".
func serviceHasReadyPods(clientset *k8s.Clientset, svc *corev1.Service) (bool, string, error) {
	podList, err := servicePods(clientset, svc)
	if err != nil {
		return false, "error", err
	}
	for _, pod := range podList.Items {
		phase := pod.Status.Phase
		if phase != "Running" {
			return false, fmt.Sprintf("pod %s is in phase %+v", pod.Name, phase), nil
		}
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if !containerStatus.Ready {
				return false, fmt.Sprintf("container %s in pod %s is not ready", containerStatus.Name, pod.Name), nil
			}
		}
	}
	return true, "ok", nil
}

func servicePods(clientset *k8s.Clientset, svc *corev1.Service) (*corev1.PodList, error) {
	podClient := clientset.CoreV1().Pods(svc.Namespace)
	selector := []string{}
	for key, value := range svc.Spec.Selector {
		selector = append(selector, fmt.Sprintf("%s=%s", key, value))
	}
	pods, err := podClient.List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: strings.Join(selector, ","),
		},
	)
	if err != nil {
		return nil, err
	}
	return pods.DeepCopy(), nil
}
