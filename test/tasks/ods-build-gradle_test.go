package tasks

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"

	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildGradle(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-gradle",
		[]tasktesting.Service{
			tasktesting.Bitbucket,
			tasktesting.Nexus,
			tasktesting.SonarQube,
		},
		map[string]tasktesting.TestCase{
			"task should build gradle app": {
				WorkspaceDirMapping: map[string]string{"source": "gradle-sample-app"},
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
						"docker/Dockerfile",
						"docker/app.jar",
						filepath.Join(pipelinectxt.XUnitReportsPath, "report.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.xml"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "quality-gate.json"),
					)

					logContains("No sonar-project.properties present, using default:", ctxt.CollectedLogs, t)
					logContains("Using NEXUS_URL=http://ods-test-nexus.kind:8081", ctxt.CollectedLogs, t)
				},
			},
		})
}

func logContains(wantLogMsg string, collectedLogs []byte, t *testing.T) {
	logString := string(collectedLogs)
	if !strings.Contains(logString, wantLogMsg) {
		t.Fatalf("Want:\n%s\n\nGot:\n%s", wantLogMsg, logString)
	}
}
