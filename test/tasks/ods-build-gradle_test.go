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
						filepath.Join(pipelinectxt.XUnitReportsPath, "TEST-ods.java.gradle.sample.app.AppTest.xml"),
						filepath.Join(pipelinectxt.XUnitReportsPath, "TEST-ods.java.gradle.sample.app.AppTest2.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "coverage.xml"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "quality-gate.json"),
					)

					logContains(ctxt.CollectedLogs, t,
						"--gradle-options=--no-daemon --stacktrace",
						"No sonar-project.properties present, using default:",
						"Using NEXUS_URL=http://ods-test-nexus.kind:8081",
						"Gradle 7.3.3",
						"Using GRADLE_OPTS=-Dorg.gradle.jvmargs=-Xmx512M",
						"To honour the JVM settings for this build a single-use Daemon process will be forked.",
						"Using ARTIFACTS_DIR=/workspace/source/.ods/artifacts",
					)
				},
			},
		})
}

func logContains(collectedLogs []byte, t *testing.T, wantLogMsgs ...string) {
	logString := string(collectedLogs)

	for _, msg := range wantLogMsgs {
		if !strings.Contains(logString, msg) {
			t.Fatalf("Want:\n%s\n\nGot:\n%s", msg, logString)
		}
	}

}
