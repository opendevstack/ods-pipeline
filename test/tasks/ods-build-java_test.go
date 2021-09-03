package tasks

import (
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildJava(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-java",
		map[string]tasktesting.TestCase{
			"task should build java gradle app": {
				WorkspaceDirMapping: map[string]string{"source": "java-gradle-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"sonar-quality-gate": "true",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					wantFiles := []string{
						"docker/Dockerfile",
						"docker/app.jar",
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.xml"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
					}
					for _, wf := range wantFiles {
						if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
							t.Fatalf("Want %s, but got nothing", wf)
						}
					}

					b, _, err := command.Run(wsDir+"/docker/app.jar", []string{})
					if err != nil {
						t.Fatal(err)
					}
					if string(b) != "Hello World" {
						t.Fatalf("Got: %+v, want: %+v.", string(b), "Hello World")
					}
				},
			},
			//"task should fail linting go app and generate lint report": {
			//	WorkspaceDirMapping: map[string]string{"source": "go-sample-app-lint-error"},
			//	PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
			//		wsDir := ctxt.Workspaces["source"]
			//		ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
			//		ctxt.Params = map[string]string{
			//			"java-image":  "localhost:5000/ods/ods-java-toolset:latest",
			//			"sonar-image": "localhost:5000/ods/ods-sonar:latest",
			//		}
			//	},
			//	WantRunSuccess: false,
			//	PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
			//		wsDir := ctxt.Workspaces["source"]
			//
			//		wantFiles := []string{
			//			".ods/artifacts/lint-report/report.txt",
			//		}
			//
			//		for _, wf := range wantFiles {
			//			if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
			//				t.Fatalf("Want %s, but got nothing", wf)
			//			}
			//		}
			//
			//		wantLintReportContent := "main.go:6:2: printf: Printf format %s reads arg #1, but call has 0 args (govet)\n\tfmt.Printf(\"Hello World %s\") // lint error on purpose to generate lint report\n\t^"
			//
			//		checkFileContent(t, wsDir, ".ods/artifacts/lint-report/report.txt", wantLintReportContent)
			//	},
			//},
		})
}
