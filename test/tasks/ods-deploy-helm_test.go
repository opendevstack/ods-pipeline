package tasks

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	"sigs.k8s.io/yaml"
)

func TestTaskODSDeployHelm(t *testing.T) {
	runTaskTestCases(t,
		"ods-deploy-helm-v0-1-0",
		map[string]tasktesting.TestCase{
			"should upgrade Helm chart": {
				WorkspaceDirMapping: map[string]string{"source": "helm-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"image": "localhost:5000/ods/ods-helm:latest",
					}

					err := createODSYML(wsDir, ctxt.Namespace)
					if err != nil {
						t.Fatal(err)
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					//wsDir := ctxt.Workspaces["source"]
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

func createODSYML(wsDir, releaseNamespace string) error {
	o := &config.ODS{
		Environments: config.Environments{
			DEV: config.Environment{
				Targets: []config.Target{
					{
						Name:      "default",
						Namespace: releaseNamespace,
					},
				},
			},
		},
	}
	y, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	filename := filepath.Join(wsDir, "ods.yml")
	return ioutil.WriteFile(filename, y, 0644)
}
