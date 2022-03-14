package tasks

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildGo(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-go",
		[]tasktesting.Service{
			tasktesting.Nexus,
			tasktesting.SonarQube,
		},
		map[string]tasktesting.TestCase{
			"build go app": {
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

					checkFilesExist(t, wsDir,
						"docker/Dockerfile",
						"docker/app",
						filepath.Join(pipelinectxt.LintReportsPath, "report.txt"),
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.out"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "quality-gate.json"),
					)

					sonarProject := sonar.ProjectKey(ctxt.ODS, "")
					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, sonarProject, true, "OK")

					wantLogMsg := "No sonar-project.properties present, using default:"
					if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
						t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
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
			"build go app in subdirectory": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					// Setup subdir in "monorepo"
					subdir := "go-src"
					createAppInSubDirectory(t, wsDir, subdir, "go-sample-app")

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

					checkFilesExist(t, wsDir,
						fmt.Sprintf("%s/docker/Dockerfile", subdir),
						binary,
						filepath.Join(pipelinectxt.LintReportsPath, fmt.Sprintf("%s-report.txt", subdir)),
						filepath.Join(pipelinectxt.XUnitReportsPath, fmt.Sprintf("%s-report.xml", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-coverage.out", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-analysis-report.md", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-issues-report.csv", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-quality-gate.json", subdir)))

					sonarProject := sonar.ProjectKey(ctxt.ODS, subdir+"-")
					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, sonarProject, true, "OK")

					b, _, err := command.Run(filepath.Join(wsDir, binary), []string{})
					if err != nil {
						t.Fatal(err)
					}
					if string(b) != "Hello World" {
						t.Fatalf("Got: %+v, want: %+v.", string(b), "Hello World")
					}
				},
			},
			"fail linting go app and generate lint report": {
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

					wantFile := filepath.Join(pipelinectxt.LintReportsPath, "report.txt")
					checkFilesExist(t, wsDir, wantFile)

					wantLintReportContent := "main.go:6:2: printf: fmt.Printf format %s reads arg #1, but call has 0 args (govet)\n\tfmt.Printf(\"Hello World %s\") // lint error on purpose to generate lint report\n\t^"

					checkFileContent(t, wsDir, ".ods/artifacts/lint-reports/report.txt", wantLintReportContent)
				},
			},
			"build go app with pre-test script": {
				WorkspaceDirMapping: map[string]string{"source": "go-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"sonar-skip":      "true",
						"pre-test-script": "pre-test-script.sh",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					wantFile := "docker/test.txt"
					checkFilesExist(t, wsDir, wantFile)
				},
			},
			"build go app in PR": {
				WorkspaceDirMapping: map[string]string{"source": "go-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					writeContextFile(t, wsDir, "pr-key", "3")
					writeContextFile(t, wsDir, "pr-key", "master")
					ctxt.Params = map[string]string{
						"go-os":              runtime.GOOS,
						"go-arch":            runtime.GOARCH,
						"sonar-quality-gate": "true",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					sonarProject := sonar.ProjectKey(ctxt.ODS, "")
					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, sonarProject, true, "OK")
				},
			},
		})
}
