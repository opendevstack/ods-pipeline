package tasks

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildNPM(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-npm",
		requiredServices(tasktesting.Nexus, tasktesting.SonarQube),
		map[string]tasktesting.TestCase{
			"build typescript app with SQ scan": {
				WorkspaceDirMapping: map[string]string{"source": "typescript-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = buildTaskParams(map[string]string{
						"sonar-quality-gate": "true",
						"cache-build":        "false",
					})
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					checkFilesExist(t, wsDir,
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "clover.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage-final.json"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "lcov.info"),
						filepath.Join(pipelinectxt.LintReportsPath, "report.txt"),
						"docker/dist/src/index.js",
						"docker/dist/package.json",
						"docker/dist/package-lock.json",
					)
					if !*skipSonarQubeFlag {
						checkFilesExist(t, wsDir,
							filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
							filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
							filepath.Join(pipelinectxt.SonarAnalysisPath, "quality-gate.json"),
						)
					}

					wantLogMsg := "No sonar-project.properties present, using default:"
					if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
						t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
					}

					if !*skipSonarQubeFlag {
						sonarProject := sonar.ProjectKey(ctxt.ODS, "")
						checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, sonarProject, true, "OK")
					}
				},
			},
			"build javascript app in subdirectory with build caching": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					// Setup subdir in "monorepo"
					subdir := "js-src"
					createAppInSubDirectory(t, wsDir, subdir, "javascript-sample-app")

					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = buildTaskParams(map[string]string{
						"working-dir": subdir,
						"output-dir":  "../docker",
					})
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					subdir := "js-src"
					checkFilesExist(t, wsDir,
						filepath.Join(pipelinectxt.XUnitReportsPath, fmt.Sprintf("%s-report.xml", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-clover.xml", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-coverage-final.json", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-lcov.info", subdir)),
						filepath.Join(pipelinectxt.LintReportsPath, fmt.Sprintf("%s-report.txt", subdir)),
						"docker/dist/src/index.js",
						"docker/dist/package.json",
						"docker/dist/package-lock.json",
					)
				},
				AdditionalRuns: []tasktesting.TaskRunCase{{
					// inherits funcs from primary task only set explicitly
					PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
						// ctxt still in place from prior run
						wsDir := ctxt.Workspaces["source"]
						tasktesting.RemoveAll(t, wsDir, "docker/dist")
						tasktesting.RemoveAll(t, wsDir, "dist")
					},
					WantRunSuccess: true,
				}},
			},
			"fail linting typescript app and generate lint report": {
				WorkspaceDirMapping: map[string]string{"source": "typescript-sample-app-lint-error"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
				},
				WantRunSuccess: false,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					wantFile := filepath.Join(pipelinectxt.LintReportsPath, "report.txt")
					checkFilesExist(t, wsDir, wantFile)

					wantLintReportContent := "/workspace/source/src/index.ts: line 3, col 31, Warning - Unexpected any. Specify a different type. (@typescript-eslint/no-explicit-any)\n\n1 problem"
					checkFileContentContains(t, wsDir, filepath.Join(pipelinectxt.LintReportsPath, "report.txt"), wantLintReportContent)
				},
			},
			"fail pulling image if unsupported node version is specified": {
				WorkspaceDirMapping: map[string]string{"source": "javascript-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"node-version": "10",
					}
				},
				WantSetupFail: true,
			},
			"build backend javascript app": {
				Timeout:             10 * time.Minute,
				WorkspaceDirMapping: map[string]string{"source": "javascript-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = buildTaskParams(map[string]string{
						"copy-node-modules": "true",
					})
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkFilesExist(t, wsDir,
						"docker/dist/node_modules/",
						"docker/dist/package.json",
						"docker/dist/package-lock.json",
					)
				},
			},
			"build javascript app with custom build directory": {
				WorkspaceDirMapping: map[string]string{"source": "javascript-sample-app-build-dir"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = buildTaskParams(map[string]string{
						"build-dir": "build",
					})
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkFilesExist(t, wsDir,
						"docker/dist/src/index.js",
						"docker/dist/package.json",
						"docker/dist/package-lock.json",
					)
				},
			},
			"build javascript app using node16": {
				WorkspaceDirMapping: map[string]string{"source": "javascript-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"sonar-skip":   "true",
						"node-version": "16",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkFilesExist(t, wsDir,
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "clover.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage-final.json"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "lcov.info"),
						filepath.Join(pipelinectxt.LintReportsPath, "report.txt"),
						"docker/dist/src/index.js",
						"docker/dist/package.json",
						"docker/dist/package-lock.json",
					)
				},
			},
		})
}
