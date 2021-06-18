package tasks

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	"sigs.k8s.io/yaml"
)

func TestTaskODSDeployHelm(t *testing.T) {

	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard", // if using KinD, set it to "standard"
			TaskDir:          projectpath.Root + "/deploy/tasks",
			EnvironmentDir:   projectpath.Root + "/test/testdata/deploy/cd-kind",
		},
	)

	tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	defer tasktesting.TearDown(t, c, ns)

	tests := map[string]tasktesting.TestCase{
		"should upgrade Helm chart": {
			WorkspaceDirMapping: map[string]string{"source": "helm-sample-app"},
			PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				ctxt.ODS = tasktesting.SetupGitRepo(t, ns, wsDir)
				ctxt.Params = map[string]string{
					"image": "localhost:5000/ods/helm:latest",
				}

				o := &config.ODS{
					Environments: config.Environments{
						DEV: config.Environment{
							Targets: []config.Target{
								{
									Name:      "default",
									Namespace: ns, // other ns should be used
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
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             "Task",                   // could be read from task definition
				TaskName:                "ods-deploy-helm-v0-1-0", // could be read from task definition
				Clients:                 c,
				Namespace:               ns,
				Timeout:                 5 * time.Minute, // depending on  the task we may need to increase or decrease it
				AlwaysKeepTmpWorkspaces: *alwaysKeepTmpWorkspacesFlag,
			})

		})

	}
}
