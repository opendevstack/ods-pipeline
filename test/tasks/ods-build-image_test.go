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
				checkResultingFiles(t, wsDir)
				checkResultingImage(t, ns, wsDir)
			},
		},
		"task should reuse existing image": {
			WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
			Params: map[string]string{
				"registry":      "kind-registry.kind:5000",
				"builder-image": "localhost:5000/ods/buildah:latest",
				"tls-verify":    "false",
			},
			WantSuccess: true,
			PrepareFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]
				buildAndPushImage(t, ns, wsDir)
			},
			CheckFunc: func(t *testing.T, workspaces map[string]string) {
				wsDir := workspaces["source"]
				checkResultingFiles(t, wsDir)
				checkResultingImage(t, ns, wsDir)
				// TODO: actually check that we did not rebuild the image ...
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

func buildAndPushImage(t *testing.T, ns, wsDir string) {
	tag := getDockerImageTag(t, ns, wsDir)
	_, stderr, err := command.Run("docker", []string{
		"build", "-t", tag, filepath.Join(wsDir, "docker"),
	})
	if err != nil {
		t.Fatalf("could not build image: %s, stderr: %s", err, string(stderr))
	}
	_, stderr, err = command.Run("docker", []string{
		"push", tag,
	})
	if err != nil {
		t.Fatalf("could not push image: %s, stderr: %s", err, string(stderr))
	}
}

func checkResultingFiles(t *testing.T, wsDir string) {
	wantFiles := []string{
		".ods/artifacts/image-digests/hello-world",
	}
	for _, wf := range wantFiles {
		if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
			t.Fatalf("Want %s, but got nothing", wf)
		}
	}
}

func checkResultingImage(t *testing.T, ns, wsDir string) {
	stdout, stderr, err := command.Run("docker", []string{
		"run", "--rm",
		getDockerImageTag(t, ns, wsDir),
	})
	if err != nil {
		t.Fatalf("could not run built image: %s, stderr: %s", err, string(stderr))
	}
	got := strings.TrimSpace(string(stdout))
	want := "Hello World"
	if got != want {
		t.Fatalf("Want %s, but got %s", want, got)
	}
}

func getDockerImageTag(t *testing.T, ns, wsDir string) string {
	sha, err := getTrimmedFileContent(filepath.Join(wsDir, ".ods/git-commit"))
	if err != nil {
		t.Fatalf("could not read git-commit: %s", err)
	}
	return fmt.Sprintf("localhost:5000/%s/hello-world:%s", ns, sha)
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
