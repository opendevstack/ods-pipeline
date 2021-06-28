package tasks

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildGo(t *testing.T) {

	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard", // if using KinD, set it to "standard"
		},
	)

	tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	defer tasktesting.TearDown(t, c, ns)

	tests := map[string]tasktesting.TestCase{
		"task should build go app": {
			WorkspaceDirMapping: map[string]string{"source": "go-sample-app"},
			PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				ctxt.ODS = tasktesting.SetupGitRepo(t, ns, wsDir)
				ctxt.Params = map[string]string{
					"go-image":    "localhost:5000/ods/ods-go-toolset:latest",
					"sonar-image": "localhost:5000/ods/ods-sonar:latest",
					"go-os":       runtime.GOOS,
					"go-arch":     runtime.GOARCH,
				}
			},
			WantRunSuccess: true,
			PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]

				wantFiles := []string{
					"docker/Dockerfile",
					"docker/app",
					"build/test-results/test/report.xml",
					"coverage.out",
					"test-results.txt",
					".ods/artifacts/xunit-reports/report.xml",
					".ods/artifacts/code-coverage/coverage.out",
					".ods/artifacts/sonarqube-analysis/analysis-report.md",
					".ods/artifacts/sonarqube-analysis/issues-report.csv",
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
		"task should fail linting go app and generate lint report": {
			WorkspaceDirMapping: map[string]string{"source": "go-sample-app-lint-error"},
			PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				ctxt.ODS = tasktesting.SetupGitRepo(t, ns, wsDir)
				ctxt.Params = map[string]string{
					"go-image":    "localhost:5000/ods/ods-go-toolset:latest",
					"sonar-image": "localhost:5000/ods/ods-sonar:latest",
					"go-os":       runtime.GOOS,
					"go-arch":     runtime.GOARCH,
				}
			},
			WantRunSuccess: false,
			PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]

				wantFiles := []string{
					".ods/artifacts/lint-report/report.txt",
				}

				for _, wf := range wantFiles {
					if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
						t.Fatalf("Want %s, but got nothing", wf)
					}
				}

				wantLintReportContent := "main.go:6:2: printf: Printf format %s reads arg #1, but call has 0 args (govet)\n\tfmt.Printf(\"Hello World %s\") // lint error on purpose to generate lint report\n\t^"

				checkFileContent(t, wsDir, ".ods/artifacts/lint-report/report.txt", wantLintReportContent)
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             "ClusterTask",         // could be read from task definition
				TaskName:                "ods-build-go-v0-1-0", // could be read from task definition
				Clients:                 c,
				Namespace:               ns,
				Timeout:                 5 * time.Minute, // depending on  the task we may need to increase or decrease it
				AlwaysKeepTmpWorkspaces: *alwaysKeepTmpWorkspacesFlag,
			})

		})

	}
}
