package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildGo(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-go",
		map[string]tasktesting.TestCase{
			"should build go app": {
				WorkspaceDirMapping: map[string]string{"source": "go-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"go-os":              runtime.GOOS,
						"go-arch":            runtime.GOARCH,
						"sonar-quality-gate": "true",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					wantFiles := []string{
						"docker/Dockerfile",
						"docker/app",
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.out"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
					}
					for _, wf := range wantFiles {
						if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
							t.Fatalf("Want %s, but got nothing", wf)
						}
					}

					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt, true, "OK")

					b, _, err := command.Run(wsDir+"/docker/app", []string{})
					if err != nil {
						t.Fatal(err)
					}
					if string(b) != "Hello World" {
						t.Fatalf("Got: %+v, want: %+v.", string(b), "Hello World")
					}
				},
			},
			"should build go app in subdirectory": {
				WorkspaceDirMapping: map[string]string{"source": "go-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					subdir := "go-src"
					err := os.MkdirAll(subdir, 0755)
					if err != nil {
						t.Fatal(err)
					}
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"go-os":              runtime.GOOS,
						"go-arch":            runtime.GOARCH,
						"sonar-quality-gate": "true",
						"working-dir":        subdir,
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					subdir := "go-src"
					binary := fmt.Sprintf("%s/docker/app", subdir)
					wantFiles := []string{
						fmt.Sprintf("%s/docker/Dockerfile", subdir),
						binary,
						filepath.Join(pipelinectxt.XUnitReportsPath, fmt.Sprintf("%s-report.xml", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-coverage.out", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-analysis-report.md", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-issues-report.csv", subdir)),
					}
					for _, wf := range wantFiles {
						if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
							t.Fatalf("Want %s, but got nothing", wf)
						}
					}

					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt, true, "OK")

					b, _, err := command.Run(filepath.Join(wsDir, binary), []string{})
					if err != nil {
						t.Fatal(err)
					}
					if string(b) != "Hello World" {
						t.Fatalf("Got: %+v, want: %+v.", string(b), "Hello World")
					}
				},
			},
			"should fail linting go app and generate lint report": {
				WorkspaceDirMapping: map[string]string{"source": "go-sample-app-lint-error"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"go-os":   runtime.GOOS,
						"go-arch": runtime.GOARCH,
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

					wantLintReportContent := "main.go:6:2: printf: fmt.Printf format %s reads arg #1, but call has 0 args (govet)\n\tfmt.Printf(\"Hello World %s\") // lint error on purpose to generate lint report\n\t^"

					checkFileContent(t, wsDir, ".ods/artifacts/lint-report/report.txt", wantLintReportContent)
				},
			},
		})
}
