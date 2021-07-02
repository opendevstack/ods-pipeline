package tasks

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/random"
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

					// TODO: defer cleanup
					releaseNamespace := random.PseudoString()
					kubernetes.CreateNamespace(ctxt.Clients.KubernetesClientSet, releaseNamespace)

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
						t.Fatal(err)
					}
					filename := filepath.Join(wsDir, "ods.yml")
					err = ioutil.WriteFile(filename, y, 0644)
					if err != nil {
						t.Fatal(err)
					}
					//k.CreateNamespace(c.KubernetesClientSet, "michael-dev")
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					//wsDir := ctxt.Workspaces["source"]
				},
			},
		},
	)
}
