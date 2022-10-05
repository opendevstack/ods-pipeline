package tasks

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"

	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildSbt(t *testing.T) {
	runTaskTestCases(t,
		"ods-build-sbt",
		[]tasktesting.Service{
			tasktesting.Nexus,
			tasktesting.SonarQube,
		},
		map[string]tasktesting.TestCase{
			"task should build sbt sample app": {
				Timeout:             10 * time.Minute,
				WorkspaceDirMapping: map[string]string{"source": "sbt-sample-app"},
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
						"docker/dist",
						filepath.Join(pipelinectxt.XUnitReportsPath, "TEST-controllers.HomeControllerSpec.xml"),
						filepath.Join(pipelinectxt.CodeCoveragesPath, "scoverage.xml"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "analysis-report.md"),
						filepath.Join(pipelinectxt.SonarAnalysisPath, "issues-report.csv"),
					)
				},
			},
		})
}
