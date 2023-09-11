package e2e

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	ott "github.com/opendevstack/ods-pipeline/pkg/odstasktest"
	"github.com/opendevstack/ods-pipeline/pkg/tasktesting"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	namespaceConfig        *ttr.NamespaceConfig
	rootPath               = "../.."
	testdataWorkspacesPath = "testdata/workspaces"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	cc, err := ttr.StartKinDCluster(
		ttr.LoadImage(ttr.ImageBuildConfig{
			Dockerfile: "build/images/Dockerfile.start",
			ContextDir: rootPath,
		}),
		ttr.LoadImage(ttr.ImageBuildConfig{
			Dockerfile: "build/images/Dockerfile.finish",
			ContextDir: rootPath,
		}),
		ttr.LoadImage(ttr.ImageBuildConfig{
			Dockerfile: "build/images/Dockerfile.pipeline-manager",
			ContextDir: rootPath,
		}),
	)
	if err != nil {
		log.Fatal("Could not start KinD cluster: ", err)
	}
	nc, cleanup, err := ttr.SetupTempNamespace(
		cc,
		ott.StartBitbucket(),
		ott.StartNexus(),
		ott.InstallODSPipeline(),
	)
	if err != nil {
		log.Fatal("Could not setup temporary namespace: ", err)
	}
	defer cleanup()
	namespaceConfig = nc
	return m.Run()
}

func newK8sClient(t *testing.T) *kubernetes.Clientset {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	kubernetesClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	return kubernetesClientset
}

func newTektonClient(t *testing.T) *tekton.Clientset {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	tektonClientSet, err := tekton.NewForConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	return tektonClientSet
}

// initBitbucketRepo initialises a Git repository inside the given workspace,
// then commits and pushes to Bitbucket.
// The workspace will also be setup with an ODS context directory in .ods
// with the given namespace.
func initBitbucketRepo(t *testing.T, k8sClient kubernetes.Interface, namespace string) ttr.WorkspaceOpt {
	return func(c *ttr.WorkspaceConfig) error {
		_ = tasktesting.SetupBitbucketRepo(t, k8sClient, namespace, c.Dir, tasktesting.BitbucketProjectKey, false)
		return nil
	}
}

// withBitbucketSourceWorkspace configures the task run with a workspace named
// "source", mapped to the directory sourced from sourceDir. The directory is
// initialised as a Git repository with an ODS context with the given namespace.
func withBitbucketSourceWorkspace(t *testing.T, sourceDir string, k8sClient kubernetes.Interface, namespace string, opts ...ttr.WorkspaceOpt) ttr.TaskRunOpt {
	return ott.WithSourceWorkspace(
		t, sourceDir,
		append([]ttr.WorkspaceOpt{initBitbucketRepo(t, k8sClient, namespace)}, opts...)...,
	)
}

func checkBuildStatus(t *testing.T, c *bitbucket.Client, gitCommit, wantBuildStatus string) {
	buildStatusPage, err := c.BuildStatusList(gitCommit)
	buildStatus := buildStatusPage.Values[0]
	if err != nil {
		t.Fatal(err)
	}
	if buildStatus.State != wantBuildStatus {
		t.Fatalf("Got: %s, want: %s", buildStatus.State, wantBuildStatus)
	}
}
