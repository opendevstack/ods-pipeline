package tasks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildTypescript(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-typescript",
		map[string]tasktesting.TestCase{
			"task should build typescript app": {
				WorkspaceDirMapping: map[string]string{"source": "typescript-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					wantFiles := []string{
						"build/test-results/test/report.xml",
						"build/coverage/clover.xml",
						"build/coverage/coverage-final.json",
						"build/coverage/lcov.info",
						".ods/artifacts/xunit-reports/report.xml",
						".ods/artifacts/code-coverage/clover.xml",
						".ods/artifacts/code-coverage/coverage-final.json",
						".ods/artifacts/code-coverage/lcov.info",
						".ods/artifacts/sonarqube-analysis/analysis-report.md",
						".ods/artifacts/sonarqube-analysis/issues-report.csv",
					}
					for _, wf := range wantFiles {
						if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
							t.Fatalf("Want %s, but got nothing", wf)
						}
					}

				},
			},
		})
}
