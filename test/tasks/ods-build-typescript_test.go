package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildTypescript(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-typescript",
		[]tasktesting.Service{
			tasktesting.Bitbucket,
			tasktesting.Nexus,
			tasktesting.SonarQube,
		},
		map[string]tasktesting.TestCase{
			"build typescript app": {
				WorkspaceDirMapping: map[string]string{"source": "typescript-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					wantFiles := []string{
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "clover.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage-final.json"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "lcov.info"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
					}
					for _, wf := range wantFiles {
						if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
							t.Fatalf("Want %s, but got nothing", wf)
						}
					}

				},
			},
			"build typescript app in subdirectory": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					// Setup subdir in "monorepo"
					subdir := "ts-src"
					err := os.MkdirAll(filepath.Join(wsDir, subdir), 0755)
					if err != nil {
						t.Fatal(err)
					}
					err = directory.Copy(
						filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "typescript-sample-app"),
						filepath.Join(wsDir, subdir),
					)
					if err != nil {
						t.Fatal(err)
					}

					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"sonar-quality-gate": "true",
						"working-dir":        subdir,
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					subdir := "ts-src"
					wantFiles := []string{
						filepath.Join(pipelinectxt.XUnitReportsPath, fmt.Sprintf("%s-report.xml", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-clover.xml", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-coverage-final.json", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-lcov.info", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-analysis-report.md", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-issues-report.csv", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-quality-gate.json", subdir)),
					}
					for _, wf := range wantFiles {
						if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
							t.Fatalf("Want %s, but got nothing", wf)
						}
					}

					sonarProject := sonar.ProjectKey(ctxt.ODS, subdir+"-")
					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, sonarProject, true, "OK")
				},
			},
		})
}
