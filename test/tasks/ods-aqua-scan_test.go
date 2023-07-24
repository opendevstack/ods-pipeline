package tasks

import (
	"testing"

	"github.com/opendevstack/ods-pipeline/pkg/tasktesting"
)

func TestTaskODSAquaScan(t *testing.T) {
	runTaskTestCases(t,
		"ods-aqua-scan",
		[]tasktesting.Service{},
		map[string]tasktesting.TestCase{
			"task fails without Aqua download URL": {
				WorkspaceDirMapping: map[string]string{"source": "empty"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
				},
				WantRunSuccess: false,
			},
		},
	)
}
