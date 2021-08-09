package tasks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildPython(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-python",
		map[string]tasktesting.TestCase{
			"task should build python flask app": {
				WorkspaceDirMapping: map[string]string{"source": "python-flask-sample-app"},
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

					checkFileContentContains(t, wsDir, "build/test-results/coverage/coverage.xml", wantContains)
					checkSonarQualityGate(t, ctxt.Clients.KubernetesClientSet, ctxt, true, "OK")
				},
			},
		})
}
