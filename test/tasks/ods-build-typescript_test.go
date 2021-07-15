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

				},
			},
		})
}
