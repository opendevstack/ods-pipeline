package tasks

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildGo(t *testing.T) {

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
		"task should build go app": {
			WorkspaceDirMapping: map[string]string{"source": "go-sample-app"},
			PrepareFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]
				os.Chdir(wsDir)
				err := tasktesting.InitAndCommitOrFatal(t, wsDir)
				if err != nil {
					t.Fatal(err)
				}

				projectKey := "FOO"
				// wsName := filepath.Base(wsDir)
				// err = tasktesting.PushToBitbucket(c.KubernetesClientSet, ns, projectKey, wsName)
				// if err != nil {
				// 	t.Fatal(err)
				// }

				err = tasktesting.WriteDotOds(wsDir, projectKey)
				if err != nil {
					t.Fatal(err)
				}
			},
			Params: map[string]string{
				//"sonar-skip":  "true",
				"go-image":    "localhost:5000/ods/go-toolset:latest",
				"sonar-image": "localhost:5000/ods/sonar:latest",
				"go-os":       runtime.GOOS,
				"go-arch":     runtime.GOARCH,
			},
			WantSuccess: true,
			CheckFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]

				wantFiles := []string{
					"docker/Dockerfile",
					"docker/app",
					"build/test-results/test/report.xml",
					"coverage.out",
					"test-results.txt",
					".ods/artifacts/xunit-reports/report.xml",
				}
				for _, wf := range wantFiles {
					if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
						t.Fatalf("Want %s, but got nothing", wf)
					}
				}

				b, _, err := command.Run(wsDir+"/docker/app", []string{})
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != "Hello World" {
					t.Fatalf("Got: %+v, want: %+v.", string(b), "Hello World")
				}
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef: "ClusterTask",       // could be read from task definition
				TaskName:    "ods-build-go-v0-1", // could be read from task definition
				Clients:     c,
				Namespace:   ns,
				Timeout:     5 * time.Minute, // depending on  the task we may need to increase or decrease it
			})

		})

	}
}
