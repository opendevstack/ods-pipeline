package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSBuildImage(t *testing.T) {

	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        "/files", // this is the dir *within* the KinD container that mounts to ${ODS_PIPELINE_DIR}/test
			StorageCapacity:  "1Gi",
			StorageClassName: "standard", // if using KinD, set it to "standard"
			TaskDir:          projectpath.Root + "/deploy/tasks",
			EnvironmentDir:   projectpath.Root + "/test/testdata/deploy/cd-kind",
		},
	)

	tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	defer tasktesting.TearDown(t, c, ns)

	tests := map[string]tasktesting.TestCase{
		"task should build image": {
			WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
			Params: map[string]string{
				"registry":      "kind-registry.kind:5000",
				"builder-image": "localhost:5000/ods/buildah:latest",
				"tls-verify":    "false",
			},
			WantSuccess: true,
			CheckFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]

				wantFiles := []string{
					".ods/artifacts/image-digests/hello-world.json",
				}
				for _, wf := range wantFiles {
					if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
						t.Fatalf("Want %s, but got nothing", wf)
					}
				}

				sha, err := getTrimmedFileContent(filepath.Join(wsDir, ".ods/git-commit"))
				if err != nil {
					t.Fatal(err)
				}
				stdout, stderr, err := command.Run("docker", []string{
					"run", "--rm",
					fmt.Sprintf("localhost:5000/%s/hello-world:%s", ns, sha),
				})
				if err != nil {
					t.Fatalf("could not run resultind docker container: %s, stderr: %s", err, string(stderr))
				}
				got := strings.TrimSpace(string(stdout))
				want := "Hello World"
				if got != want {
					t.Fatalf("Want %s, but got %s", want, got)
				}
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef: "ClusterTask",          // could be read from task definition
				TaskName:    "ods-build-image-v0-1", // could be read from task definition
				Clients:     c,
				Namespace:   ns,
				Timeout:     5 * time.Minute, // depending on  the task we may need to increase or decrease it
			})

		})

	}
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
