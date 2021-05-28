package framework

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"testing"

	k "github.com/opendevstack/pipeline/internal/kubernetes"
	uuid "github.com/satori/go.uuid"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/test/logging"
)

type SetupOpts struct {
	SourceDir        string
	StorageCapacity  string
	StorageClassName string
	TaskDir          string
	EnvironmentDir   string
}

func Setup(t *testing.T, opts SetupOpts) (*k.Clients, string) {
	t.Helper()

	namespace := uuid.NewV4().String()
	clients := k.NewClients()

	k.CreateNamespace(clients.KubernetesClientSet, namespace)

	_, err := k.CreatePersistentVolume(clients.KubernetesClientSet, namespace, opts.StorageCapacity, opts.SourceDir, opts.StorageClassName)
	if err != nil {
		t.Error(err)
	}

	_, err = k.CreatePersistentVolumeClaim(clients.KubernetesClientSet, opts.StorageCapacity, opts.StorageClassName, namespace)
	if err != nil {
		t.Error(err)
	}

	applyYAMLFilesInDir(t, namespace, opts.TaskDir)

	if len(opts.EnvironmentDir) > 0 {
		applyYAMLFilesInDir(t, namespace, opts.EnvironmentDir)
	}

	return clients, namespace
}

func applyYAMLFilesInDir(t *testing.T, ns string, fileDir string) {

	stdout, stderr, err := runCmd("kubectl", []string{"-n", ns, "apply", "-f", fileDir})

	t.Logf(string(stdout))
	t.Logf(string(stderr))
	if err != nil {
		t.Error(err)
	}
}

func runCmd(executable string, args []string) (outBytes, errBytes []byte, err error) {
	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outBytes = stdout.Bytes()
	errBytes = stderr.Bytes()
	return outBytes, errBytes, err
}

func Header(logf logging.FormatLogger, text string) {
	left := "### "
	right := " ###"
	txt := left + text + right
	bar := strings.Repeat("#", len(txt))
	logf(bar)
	logf(txt)
	logf(bar)
}

// CleanupOnInterrupt will execute the function cleanup if an interrupt signal is caught
func CleanupOnInterrupt(cleanup func(), logf logging.FormatLogger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logf("Test interrupted, cleaning up.")
			cleanup()
			os.Exit(1)
		}
	}()
}

func TearDown(t *testing.T, cs *k.Clients, namespace string) {
	t.Helper()
	if cs.KubernetesClientSet == nil {
		return
	}

	t.Logf("Deleting namespace %s", namespace)
	if err := cs.KubernetesClientSet.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{}); err != nil {
		t.Errorf("Failed to delete namespace %s: %s", namespace, err)
	}

	// For simplicity and traceability, we use for the PV the same name as the namespace
	pvName := namespace
	t.Logf("Deleting persistent volume with name %s", pvName)
	if err := cs.KubernetesClientSet.CoreV1().PersistentVolumes().Delete(context.Background(), pvName, metav1.DeleteOptions{}); err != nil {
		t.Errorf("Failed to delete persistent volume %s: %s", pvName, err)
	}

}

func getCRDYaml(cs *k.Clients, ns string) ([]byte, error) {
	var output []byte
	printOrAdd := func(kind, name string, i interface{}) {
		bs, err := yaml.Marshal(i)
		if err != nil {
			return
		}
		output = append(output, []byte("\n---\n")...)
		output = append(output, bs...)
	}

	ps, err := cs.TektonClientSet.TektonV1beta1().Pipelines(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get pipeline: %w", err)
	}
	for _, i := range ps.Items {
		printOrAdd("Pipeline", i.Name, i)
	}

	prrs, err := cs.TektonClientSet.TektonV1beta1().PipelineRuns(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get pipelinerun: %w", err)
	}
	for _, i := range prrs.Items {
		printOrAdd("PipelineRun", i.Name, i)
	}

	cts, err := cs.TektonClientSet.TektonV1beta1().ClusterTasks().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get cluster tasks: %w", err)
	}
	for _, i := range cts.Items {
		printOrAdd("Task", i.Name, i)
	}

	ts, err := cs.TektonClientSet.TektonV1beta1().Tasks(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get tasks: %w", err)
	}
	for _, i := range ts.Items {
		printOrAdd("Task", i.Name, i)
	}

	trs, err := cs.TektonClientSet.TektonV1beta1().TaskRuns(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get taskrun: %w", err)
	}
	for _, i := range trs.Items {
		printOrAdd("TaskRun", i.Name, i)
	}

	pods, err := cs.KubernetesClientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get pods: %w", err)
	}
	for _, i := range pods.Items {
		printOrAdd("Pod", i.Name, i)
	}

	return output, nil
}

func CollectTaskResultInfo(tr *v1beta1.TaskRun, logf logging.FormatLogger) {
	logf("Status: %s\n", tr.Status.GetCondition(apis.ConditionSucceeded).Status)
	logf("Reason: %s\n", tr.Status.GetCondition(apis.ConditionSucceeded).GetReason())
	logf("Message: %s\n", tr.Status.GetCondition(apis.ConditionSucceeded).GetMessage())
}
