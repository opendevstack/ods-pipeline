package tasks

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

var alwaysKeepTmpWorkspacesFlag = flag.Bool("always-keep-tmp-workspaces", false, "Whether to keep temporary workspaces from taskruns even when test is successful")

const (
	bitbucketProjectKey = "ODSPIPELINETEST"
	taskKindRef         = "ClusterTask"
	storageClasName     = "standard" // if using KinD, set it to "standard"
	storageCapacity     = "1Gi"
	storageSourceDir    = "/files" // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
)

func checkFileContent(t *testing.T, wsDir, filename, want string) {
	got, err := getTrimmedFileContent(filepath.Join(wsDir, filename))
	if err != nil {
		t.Fatalf("could not read %s: %s", filename, err)
	}
	if got != want {
		t.Fatalf("got '%s', want '%s' in file %s", got, want, filename)
	}
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func runTaskTestCases(t *testing.T, taskName string, testCases map[string]tasktesting.TestCase) {
	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        storageSourceDir,
			StorageCapacity:  storageCapacity,
			StorageClassName: storageClasName,
		},
	)

	tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	defer tasktesting.TearDown(t, c, ns)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             taskKindRef, // could be read from task definition
				TaskName:                taskName,    // could be read from task definition
				Clients:                 c,
				Namespace:               ns,
				Timeout:                 5 * time.Minute, // depending on  the task we may need to increase or decrease it
				AlwaysKeepTmpWorkspaces: *alwaysKeepTmpWorkspacesFlag,
			})
		})
	}
}
