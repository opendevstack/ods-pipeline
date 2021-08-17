package tasktesting

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/internal/random"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/test/logging"
)

type SetupOpts struct {
	SourceDir        string
	StorageCapacity  string
	StorageClassName string
}

func Setup(t *testing.T, opts SetupOpts) (*k.Clients, string) {
	t.Helper()

	namespace := random.PseudoString()
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

	installCDNamespaceResources(
		t, namespace, "pipeline", "./chart/values.kind.yaml,./chart/values.generated.yaml",
	)

	return clients, namespace
}

func installCDNamespaceResources(t *testing.T, ns, serviceaccount, valuesFile string) {

	scriptArgs := []string{"-n", ns, "-s", serviceaccount, "-f", valuesFile, "--no-diff"}
	if testing.Verbose() {
		scriptArgs = append(scriptArgs, "-v")
	}

	stdout, stderr, err := command.Run(
		filepath.Join(projectpath.Root, "scripts/install-cd-namespace-resources.sh"),
		scriptArgs,
	)

	t.Logf(string(stdout))
	if err != nil {
		t.Logf(string(stderr))
		t.Fatal(err)
	}
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

func CollectTaskResultInfo(tr *v1beta1.TaskRun, logf logging.FormatLogger) {
	if tr == nil {
		logf("error: no taskrun")
		return
	}
	logf("Status: %s\n", tr.Status.GetCondition(apis.ConditionSucceeded).Status)
	logf("Reason: %s\n", tr.Status.GetCondition(apis.ConditionSucceeded).GetReason())
	logf("Message: %s\n", tr.Status.GetCondition(apis.ConditionSucceeded).GetMessage())
}
