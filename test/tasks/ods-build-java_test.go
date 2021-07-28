package tasks

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildJava(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-java-v0-1-0",
		map[string]tasktesting.TestCase{
			"task should build java/maven app": {
				WorkspaceDirMapping: map[string]string{"source": "java-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"java-image":    "localhost:5000/ods/ods-java-toolset:latest",
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
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"java-image": "localhost:5000/ods/ods-java-toolset:latest",
						"sonar-image": "localhost:5000/ods/ods-sonar:latest",
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
		})
}
