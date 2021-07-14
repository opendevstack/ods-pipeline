package tasks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					wantFiles := []string{
						"docker/app/main.py",
						"docker/app/requirements.txt",
						"build/test-results/test/report.xml",
						"build/test-results/coverage/coverage.xml",
						".ods/artifacts/xunit-reports/report.xml",
						".ods/artifacts/code-coverage/coverage.xml",
						".ods/artifacts/sonarqube-analysis/analysis-report.md",
						".ods/artifacts/sonarqube-analysis/issues-report.csv",
					}
					for _, wf := range wantFiles {
						if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
							t.Fatalf("Want %s, but got nothing", wf)
						}
					}

					wantContains := `
					<package name="." line-rate="0.8571" branch-rate="0.5" complexity="0">
						<classes>
							<class name="main.py" filename="main.py" complexity="0" line-rate="0.8571" branch-rate="0.5">
								<methods/>
								<lines>
										<line number="2" hits="1"/>
										<line number="4" hits="1"/>
										<line number="7" hits="1"/>
										<line number="8" hits="1"/>
										<line number="9" hits="1"/>
										<line number="13" hits="1" branch="true" condition-coverage="50% (1/2)" missing-branches="14"/>
										<line number="14" hits="0"/>
								</lines>
							</class>
						</classes>
					</package>`

					wantContains = strings.ReplaceAll(wantContains, "\t", "")
					wantContains = strings.ReplaceAll(wantContains, "\n", "")

					checkFileContentContains(t, wsDir, "build/test-results/coverage/coverage.xml", wantContains)
				},
			},
		})
}
