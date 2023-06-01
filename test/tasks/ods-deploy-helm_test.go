package tasks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/internal/random"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	localRegistry = "localhost:5000"
	kindRegistry  = "kind-registry.kind:5000"
)

type imageImportParams struct {
	externalRef string
	namespace   string
	workdir     string
}

func TestTaskODSDeployHelm(t *testing.T) {
	var separateReleaseNamespace string
	runTaskTestCases(t,
		"ods-deploy-helm",
		[]tasktesting.Service{},
		map[string]tasktesting.TestCase{
			"skips when no namespace is given": {
				WorkspaceDirMapping: map[string]string{"source": "helm-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					// no "namespace" param set
				},
				WantRunSuccess: true,
			},
			"upgrades Helm chart in separate namespace": {
				WorkspaceDirMapping: map[string]string{"source": "helm-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)

					externalNamespace, cleanupFunc := createReleaseNamespaceOrFatal(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace,
					)
					separateReleaseNamespace = externalNamespace
					ctxt.Cleanup = cleanupFunc
					ctxt.Params = map[string]string{
						"namespace": externalNamespace,
					}
					importImage(t, imageImportParams{
						externalRef: "index.docker.io/crccheck/hello-world",
						namespace:   ctxt.Namespace,
						workdir:     wsDir,
					})
					createSampleAppPrivateKeySecret(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkFileContentContains(
						t, wsDir,
						filepath.Join(pipelinectxt.DeploymentsPath, fmt.Sprintf("diff-%s.txt", separateReleaseNamespace)),
						"Release was not present in Helm.  Diff will show entire contents as new.",
						"Deployment (apps) has been added",
						"Secret (v1) has been added",
						"Service (v1) has been added",
					)
					checkFileContentContains(
						t, wsDir,
						filepath.Join(pipelinectxt.DeploymentsPath, fmt.Sprintf("release-%s.txt", separateReleaseNamespace)),
						"Installing it now.",
						fmt.Sprintf("NAMESPACE: %s", separateReleaseNamespace),
						"STATUS: deployed",
						"REVISION: 1",
					)
					resourceName := fmt.Sprintf("%s-%s", ctxt.ODS.Component, "helm-sample-app")
					_, err := checkService(ctxt.Clients.KubernetesClientSet, separateReleaseNamespace, resourceName)
					if err != nil {
						t.Fatal(err)
					}
					_, err = checkDeployment(ctxt.Clients.KubernetesClientSet, separateReleaseNamespace, resourceName)
					if err != nil {
						t.Fatal(err)
					}

					// Verify log output massaging
					doNotWantLogMsg := "plugin \"diff\" exited with error"
					if strings.Contains(string(ctxt.CollectedLogs), doNotWantLogMsg) {
						t.Fatalf("Do not want:\n%s\n\nGot:\n%s", doNotWantLogMsg, string(ctxt.CollectedLogs))
					}
					wantLogMsg := "identified at least one change"
					if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
						t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
					}
				},
			},
			"upgrades Helm chart with dependencies": {
				WorkspaceDirMapping: map[string]string{"source": "helm-app-with-dependencies"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"namespace": ctxt.Namespace,
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					parentChartResourceName := fmt.Sprintf("%s-%s", ctxt.ODS.Component, "helm-app-with-dependencies")
					// Parent chart
					_, err := checkService(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, parentChartResourceName)
					if err != nil {
						t.Fatal(err)
					}
					_, err = checkDeployment(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, parentChartResourceName)
					if err != nil {
						t.Fatal(err)
					}
					// Subchart
					subChartResourceName := "helm-sample-database" // fixed name due to fullnameOverride
					_, err = checkService(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, subChartResourceName)
					if err != nil {
						t.Fatal(err)
					}
					d, err := checkDeployment(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, subChartResourceName)
					if err != nil {
						t.Fatal(err)
					}
					// Check that Helm value overriding in subchart works
					gotEnvValue := d.Spec.Template.Spec.Containers[0].Env[0].Value
					wantEnvValue := "tom" // defined in parent (child has value "john")
					if gotEnvValue != wantEnvValue {
						t.Fatalf("Want ENV username = %s, got: %s", wantEnvValue, gotEnvValue)
					}
				},
				AdditionalRuns: []tasktesting.TaskRunCase{{
					// inherits funcs from primary task only set explicitly
					PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
						// ctxt still in place from prior run
					},
					WantRunSuccess: true,
					PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
						wantLogMsg := "No diff detected, skipping helm upgrade"
						if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
							t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
						}
					},
				}},
			},
			"skips upgrade when diff-only is requested": {
				WorkspaceDirMapping: map[string]string{"source": "helm-app-with-dependencies"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					externalNamespace, cleanupFunc := createReleaseNamespaceOrFatal(
						t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace,
					)
					separateReleaseNamespace = externalNamespace
					ctxt.Cleanup = cleanupFunc
					ctxt.Params = map[string]string{
						"namespace": externalNamespace,
						"diff-only": "true",
					}
					importImage(t, imageImportParams{
						externalRef: "index.docker.io/crccheck/hello-world",
						namespace:   ctxt.Namespace,
						workdir:     wsDir,
					})
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					t.Log("Verify image was not promoted ...")
					img := fmt.Sprintf("%s/%s/hello-world", localRegistry, separateReleaseNamespace)
					promoted := checkIfImageExists(t, img)
					if promoted {
						t.Fatalf("Image %s should not have been promoted to %s", img, separateReleaseNamespace)
					}
					t.Log("Verify service was not deployed ...")
					resourceName := fmt.Sprintf("%s-%s", ctxt.ODS.Component, "helm-app-with-dependencies")
					_, err := checkService(ctxt.Clients.KubernetesClientSet, separateReleaseNamespace, resourceName)
					if err == nil {
						t.Fatalf("Service %s should not have been deployed to %s", resourceName, separateReleaseNamespace)
					}
					t.Log("Verify task skipped upgrade ...")
					wantLogMsg := "Only diff was requested, skipping helm upgrade"
					if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
						t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
					}
				},
			},
		},
	)
}

func createSampleAppPrivateKeySecret(t *testing.T, clientset *k8s.Clientset, ctxtNamespace string) {
	secret, err := readPrivateKeySecret()
	if err != nil {
		t.Fatal(err)
	}
	_, err = kubernetes.CreateSecret(clientset, ctxtNamespace, secret)
	if err != nil {
		t.Fatal(err)
	}
}

func createReleaseNamespaceOrFatal(t *testing.T, clientset *k8s.Clientset, ctxtNamespace string) (string, func()) {
	externalNamespace, err := createReleaseNamespace(clientset, ctxtNamespace)
	if err != nil {
		t.Fatal(err)
	}
	return externalNamespace, func() {
		if err := clientset.CoreV1().Namespaces().Delete(context.TODO(), externalNamespace, metav1.DeleteOptions{}); err != nil {
			t.Errorf("Failed to delete namespace %s: %s", externalNamespace, err)
		}
	}
}

func createReleaseNamespace(clientset *k8s.Clientset, ctxtNamespace string) (string, error) {
	releaseNamespace := random.PseudoString()
	kubernetes.CreateNamespace(clientset, releaseNamespace)
	_, err := clientset.RbacV1().RoleBindings(releaseNamespace).Create(
		context.Background(),
		&rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pipeline-deployer",
				Namespace: releaseNamespace,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "pipeline",
					Namespace: ctxtNamespace,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "edit",
			},
		},
		metav1.CreateOptions{})

	return releaseNamespace, err
}

func writeContextFile(t *testing.T, wsDir, file, content string) {
	err := os.WriteFile(
		filepath.Join(wsDir, pipelinectxt.BaseDir, file), []byte(content), 0644,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func checkDeployment(clientset *k8s.Clientset, namespace, name string) (*appsv1.Deployment, error) {
	return clientset.AppsV1().
		Deployments(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
}

func checkService(clientset *k8s.Clientset, namespace, name string) (*corev1.Service, error) {
	return clientset.CoreV1().
		Services(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
}

func readPrivateKeySecret() (*corev1.Secret, error) {
	bytes, err := os.ReadFile(filepath.Join(projectpath.Root, "test/testdata/fixtures/tasks/secret.yaml"))
	if err != nil {
		return nil, err
	}

	var secretSpec corev1.Secret
	err = yaml.Unmarshal(bytes, &secretSpec)
	if err != nil {
		return nil, err
	}
	return &secretSpec, nil
}

func importImage(t *testing.T, iip imageImportParams) {
	var err error
	cmds := [][]string{
		{"pull", iip.externalRef},
		{"tag", iip.externalRef, iip.internalRef(localRegistry)},
		{"push", iip.internalRef(localRegistry)},
	}
	for _, args := range cmds {
		if err == nil {
			_, _, err = command.RunBuffered("docker", args)
		}
	}
	if err != nil {
		t.Fatalf("docker cmd failed: %s", err)
	}

	err = pipelinectxt.WriteJsonArtifact(artifact.Image{
		Ref:        iip.internalRef(kindRegistry),
		Registry:   kindRegistry,
		Repository: iip.namespace,
		Name:       iip.name(),
		Tag:        "latest",
		Digest:     "not needed",
	}, filepath.Join(iip.workdir, pipelinectxt.ImageDigestsPath), fmt.Sprintf("%s.json", iip.name()))
	if err != nil {
		t.Fatalf("failed to write artifact: %s", err)
	}
	t.Log("Imported image", iip.internalRef(localRegistry))
}

func checkIfImageExists(t *testing.T, name string) bool {
	t.Helper()
	_, _, err := command.RunBuffered("docker", []string{"inspect", name})
	return err == nil
}

func (iip imageImportParams) name() string {
	parts := strings.Split(iip.externalRef, "/")
	return parts[2]
}

func (iip imageImportParams) internalRef(registry string) string {
	parts := strings.Split(iip.externalRef, "/")
	return fmt.Sprintf("%s/%s/%s", registry, iip.namespace, parts[2])
}
