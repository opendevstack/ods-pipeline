package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildPython(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-python",
		[]tasktesting.Service{
			tasktesting.Nexus,
			tasktesting.SonarQube,
		},
		map[string]tasktesting.TestCase{
			"build python fastapi app": {
				WorkspaceDirMapping: map[string]string{"source": "python-fastapi-sample-app"},
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

					checkFilesExist(t, wsDir,
						"docker/app/main.py",
						"docker/app/requirements.txt",
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.xml"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "quality-gate.json"),
					)

					wantContainsBytes, err := os.ReadFile("../../test/testdata/golden/ods-build-python/excerpt-from-coverage.xml")
					if err != nil {
						t.Fatal(err)
					}

					wantContains := string(wantContainsBytes)

					wantContains = strings.ReplaceAll(wantContains, "\t", "")
					wantContains = strings.ReplaceAll(wantContains, "\n", "")
					wantContains = strings.ReplaceAll(wantContains, " ", "")

					checkFileContentLeanContains(t, wsDir, filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.xml"), wantContains)
					sonarProject := sonar.ProjectKey(ctxt.ODS, "")
					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, sonarProject, true, "OK")
					wantLogMsg := "No sonar-project.properties present, using default:"
					if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
						t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
					}
				},
			},
			"build python fastapi app with build caching": {
				WorkspaceDirMapping: map[string]string{"source": "python-fastapi-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"sonar-quality-gate": "true",
						"cache-build":        "true",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					checkFilesExist(t, wsDir,
						"docker/app/main.py",
						"docker/app/requirements.txt",
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.xml"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "quality-gate.json"),
					)

					wantContainsBytes, err := os.ReadFile("../../test/testdata/golden/ods-build-python/excerpt-from-coverage.xml")
					if err != nil {
						t.Fatal(err)
					}

					wantContains := string(wantContainsBytes)

					wantContains = strings.ReplaceAll(wantContains, "\t", "")
					wantContains = strings.ReplaceAll(wantContains, "\n", "")
					wantContains = strings.ReplaceAll(wantContains, " ", "")

					checkFileContentLeanContains(t, wsDir, filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.xml"), wantContains)
					sonarProject := sonar.ProjectKey(ctxt.ODS, "")
					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, sonarProject, true, "OK")

					// This is not available when build skipping as the default is
					// supplied on the second repeat.
					// Not sure whether the check is significant in the first place.
					// wantLogMsg := "No sonar-project.properties present, using default:"
					// if !strings.Contains(string(ctxt.CollectedLogs), wantLogMsg) {
					// 	t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, string(ctxt.CollectedLogs))
					// }
				},
				AdditionalRuns: []tasktesting.TaskRunCase{{
					// inherits funcs from primary task only set explicitly
					PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
						// ctxt still in place from prior run
					},
					WantRunSuccess: true,
				}},
			},
			"build python fastapi app in subdirectory": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					// Setup subdir in "monorepo"
					subdir := "fastapi-src"
					createAppInSubDirectory(t, wsDir, subdir, "python-fastapi-sample-app")

					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"sonar-quality-gate": "true",
						"working-dir":        subdir,
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					subdir := "fastapi-src"

					checkFilesExist(t, wsDir,
						fmt.Sprintf("%s/docker/app/main.py", subdir),
						fmt.Sprintf("%s/docker/app/requirements.txt", subdir),
						filepath.Join(pipelinectxt.XUnitReportsPath, fmt.Sprintf("%s-report.xml", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-coverage.xml", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-analysis-report.md", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-issues-report.csv", subdir)),
						filepath.Join(pipelinectxt.SonarAnalysisPath, fmt.Sprintf("%s-quality-gate.json", subdir)),
					)

					sonarProject := sonar.ProjectKey(ctxt.ODS, subdir+"-")
					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace, sonarProject, true, "OK")
				},
			},
			"build python fastapi app with pre-test script": {
				WorkspaceDirMapping: map[string]string{"source": "python-fastapi-sample-app"},
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
		})
}
