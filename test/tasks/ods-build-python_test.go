package tasks

import (
	"os"
	"path/filepath"
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
						// "docker/Dockerfile",
						// "docker/app",
						// "build/test-results/test/report.xml",
						"coverage.xml",
						"test-results.txt",
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

					// TODO: Run Python Flask app
					// b, _, err := command.Run(wsDir+"/docker/app", []string{})
					// if err != nil {
					// 	t.Fatal(err)
					// }
					// if string(b) != "Hello World" {
					// 	t.Fatalf("Got: %+v, want: %+v.", string(b), "Hello World")
					// }
				},
			},
		})
}
