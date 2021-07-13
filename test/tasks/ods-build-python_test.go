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
					ctxt.Params = map[string]string{
						"no-proxy":    "127.0.0.1,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,localhost,boehringer.com,eu.boehringer.com,am.boehringer.com,biscrum.com",
						"https-proxy": "http://appaccess-zscaler.boehringer.com:80",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					wantFiles := []string{
						"docker/app/main.py",
						"docker/app/requirements.txt",
						"report.xml",
						"coverage.xml",
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
