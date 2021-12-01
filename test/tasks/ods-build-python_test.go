package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildPython(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-python",
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

					wantFiles := []string{
						"docker/app/main.py",
						"docker/app/requirements.txt",
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.xml"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "quality-gate.json"),
					}
					for _, wf := range wantFiles {
						if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
							t.Fatalf("Want %s, but got nothing", wf)
						}
					}

					wantContainsBytes, err := ioutil.ReadFile("../../test/testdata/golden/ods-build-python/excerpt-from-coverage.xml")
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
				},
			},
			"build python fastapi app in subdirectory": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					// Setup subdir in "monorepo"
					subdir := "fastapi-src"
					err := os.MkdirAll(filepath.Join(wsDir, subdir), 0755)
					if err != nil {
						t.Fatal(err)
					}
					err = directory.Copy(
						filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, "python-fastapi-sample-app"),
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
					subdir := "fastapi-src"
					wantFiles := []string{
						fmt.Sprintf("%s/docker/app/main.py", subdir),
						fmt.Sprintf("%s/docker/app/requirements.txt", subdir),
						filepath.Join(pipelinectxt.XUnitReportsPath, fmt.Sprintf("%s-report.xml", subdir)),
						filepath.Join(pipelinectxt.CodeCoveragesPath, fmt.Sprintf("%s-coverage.xml", subdir)),
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
					if _, err := os.Stat(filepath.Join(wsDir, wantFile)); os.IsNotExist(err) {
						t.Fatalf("Want %s, but got nothing", wantFile)
					}
				},
			},
		})
}
