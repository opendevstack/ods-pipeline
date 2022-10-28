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
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func TestTaskODSDeployHelm(t *testing.T) {
	var separateReleaseNamespace string
	runTaskTestCases(t,
		"ods-deploy-helm",
		[]tasktesting.Service{},
		map[string]tasktesting.TestCase{
			"skips when no environment selected": {
				WorkspaceDirMapping: map[string]string{"source": "helm-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					// simulate empty environment
					writeContextFile(t, wsDir, "environment", "")
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wantLogMsg := "Skipping deployment ..."
					if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
						t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
					}
				},
			},
			"skips when there is no diff": {
				WorkspaceDirMapping: map[string]string{"source": "helm-app-minimal"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					fatalIfErr(t, createHelmODSYML(wsDir, ctxt.Namespace))
					buildAndPushImageWithLabel(t, ctxt, "foo", wsDir)
					ia := artifact.Image{
						Image:      "http://localhost:5000/foo/bar:baz",
						Registry:   "http://localhost:5000",
						Repository: "foo",
						Name:       "bar",
						Tag:        "baz",
						Digest:     "abcdef",
					}
					fatalIfErr(t, pipelinectxt.WriteJsonArtifact(ia, pipelinectxt.ImageDigestsPath, "bar.json"))
				},
				WantRunSuccess: true,
				AdditionalRuns: []tasktesting.TaskRunCase{{
					WantRunSuccess: true,
					PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
						wantLogMsg := "Skipping deployment ..."
						if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
							t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
						}
						if imageExists(t, "a") {
							t.Fatalf("Did not expect image %s to exist in %s", "a", "b")
						}
					},
				}},
			},
			"upgrades Helm chart in separate namespace": {
				WorkspaceDirMapping: map[string]string{"source": "helm-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)

					externalNamespace, err := createReleaseNamespace(ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
					fatalIfErr(t, err)
					separateReleaseNamespace = externalNamespace
					ctxt.Cleanup = func() {
						if err := ctxt.Clients.KubernetesClientSet.CoreV1().Namespaces().Delete(context.TODO(), externalNamespace, metav1.DeleteOptions{}); err != nil {
							t.Fatalf("Failed to delete namespace %s: %s", externalNamespace, err)
						}
					}

					fatalIfErr(t, createHelmODSYML(wsDir, externalNamespace))

					createSampleAppPrivateKeySecret(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkFileContentContains(
						t, wsDir,
						filepath.Join(pipelinectxt.DeploymentsPath, "diff-dev.txt"),
						"Release was not present in Helm.  Diff will show entire contents as new.",
						"Deployment (apps) has been added",
						"Secret (v1) has been added",
						"Service (v1) has been added",
					)
					checkFileContentContains(
						t, wsDir,
						filepath.Join(pipelinectxt.DeploymentsPath, "release-dev.txt"),
						"Installing it now.",
						fmt.Sprintf("NAMESPACE: %s", separateReleaseNamespace),
						"STATUS: deployed",
						"REVISION: 1",
					)
					resourceName := fmt.Sprintf("%s-%s", ctxt.ODS.Component, "helm-sample-app")
					_, err := checkService(ctxt.Clients.KubernetesClientSet, separateReleaseNamespace, resourceName)
					fatalIfErr(t, err)
					_, err = checkDeployment(ctxt.Clients.KubernetesClientSet, separateReleaseNamespace, resourceName)
					fatalIfErr(t, err)

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
					fatalIfErr(t, createHelmODSYML(wsDir, ctxt.Namespace))
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					parentChartResourceName := fmt.Sprintf("%s-%s", ctxt.ODS.Component, "helm-app-with-dependencies")
					// Parent chart
					_, err := checkService(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, parentChartResourceName)
					fatalIfErr(t, err)
					_, err = checkDeployment(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, parentChartResourceName)
					fatalIfErr(t, err)
					// Subchart
					subChartResourceName := "helm-sample-database" // fixed name due to fullnameOverride
					_, err = checkService(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, subChartResourceName)
					fatalIfErr(t, err)
					d, err := checkDeployment(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, subChartResourceName)
					fatalIfErr(t, err)
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
	fatalIfErr(t, os.WriteFile(
		filepath.Join(wsDir, pipelinectxt.BaseDir, file), []byte(content), 0644,
	))
}

func createHelmODSYML(wsDir, releaseNamespace string) error {
	o := &config.ODS{
		Environments: []config.Environment{
			{
				Name:      "dev",
				Namespace: releaseNamespace,
				Stage:     "dev",
			},
		},
	}
	return createODSYML(wsDir, o)
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

func imageExists(t *testing.T, imageRef string) bool {
	args := []string{"inspect", "--tls-verify=false", "docker://" + imageRef}
	stdout, stderr, err := command.RunBuffered("skopeo", args)
	t.Log(string(stdout), string(stderr))
	return err == nil
}
