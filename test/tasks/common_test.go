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

func checkFileContentContains(t *testing.T, wsDir, filename, wantContains string) {
	got, err := getFileContentLean(filepath.Join(wsDir, filename))
	if err != nil {
		t.Fatalf("could not read %s: %s", filename, err)
	}
	if !strings.Contains(got, wantContains) {
		t.Fatalf("got '%s', wantContains '%s' in file %s", got, wantContains, filename)
	}
}

func getFileContentLean(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	contentStr := strings.ReplaceAll(string(content), "\t", "")
	contentStr = strings.ReplaceAll(contentStr, "\n", "")
	contentStr = strings.ReplaceAll(contentStr, " ", "")

	return contentStr, nil
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
			start := time.Now()
			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             taskKindRef, // could be read from task definition
				TaskName:                taskName,    // could be read from task definition
				Clients:                 c,
				Namespace:               ns,
				Timeout:                 5 * time.Minute, // depending on  the task we may need to increase or decrease it
				AlwaysKeepTmpWorkspaces: *alwaysKeepTmpWorkspacesFlag,
			})
			t.Logf("Test execution time: %fs", time.Since(start).Seconds())
		})
	}
}
