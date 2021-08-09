package tasks

import (
	"context"
	"fmt"
	"testing"

	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestTaskODSDeployHelm(t *testing.T) {
	runTaskTestCases(t,
		"ods-deploy-helm",
		map[string]tasktesting.TestCase{
			"should upgrade Helm chart": {
				WorkspaceDirMapping: map[string]string{"source": "helm-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)

					err := createHelmODSYML(wsDir, ctxt.Namespace)
					if err != nil {
						t.Fatal(err)
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					resourceName := fmt.Sprintf("%s-%s", ctxt.ODS.Component, "helm-sample-app")
					_, err := checkService(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, resourceName)
					if err != nil {
						t.Fatal(err)
					}
					_, err = checkDeployment(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, resourceName)
					if err != nil {
						t.Fatal(err)
					}
				},
			},
			"should upgrade Helm chart with dependencies": {
				WorkspaceDirMapping: map[string]string{"source": "helm-app-with-dependencies"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)

					err := createHelmODSYML(wsDir, ctxt.Namespace)
					if err != nil {
						t.Fatal(err)
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
					subChartResourceName := fmt.Sprintf("%s-%s", ctxt.ODS.Component, "helm-sample-database")
					_, err = checkService(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, subChartResourceName)
					if err != nil {
						t.Fatal(err)
					}
					_, err = checkDeployment(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, subChartResourceName)
					if err != nil {
						t.Fatal(err)
					}
				},
			},
		},
	)
}

// func createReleaseNamespace(clientset *kubernetes.Clientset, ctxtNamespace string) (string, error) {
// 	releaseNamespace := random.PseudoString()
// 	k.CreateNamespace(clientset, releaseNamespace)
// 	_, err := clientset.RbacV1().RoleBindings(releaseNamespace).Create(
// 		context.Background(),
// 		&v1.RoleBinding{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name:      "pipeline-deployer",
// 				Namespace: releaseNamespace,
// 			},
// 			Subjects: []v1.Subject{
// 				{
// 					Kind:      "ServiceAccount",
// 					Name:      "pipeline",
// 					Namespace: ctxtNamespace,
// 				},
// 			},
// 			RoleRef: v1.RoleRef{
// 				APIGroup: "rbac.authorization.k8s.io",
// 				Kind:     "ClusterRole",
// 				Name:     "edit",
// 			},
// 		},
// 		metav1.CreateOptions{})

// 	return releaseNamespace, err
// }

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

func checkDeployment(clientset *kubernetes.Clientset, namespace, name string) (*appsv1.Deployment, error) {
	return clientset.AppsV1().
		Deployments(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
}

func checkService(clientset *kubernetes.Clientset, namespace, name string) (*corev1.Service, error) {
	return clientset.CoreV1().
		Services(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
}
