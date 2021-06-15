package tasks

import (
	"fmt"
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
			PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				tasktesting.InitAndCommitOrFatal(t, wsDir)
				tasktesting.WriteDotOdsOrFatal(t, wsDir, bitbucketProjectKey)

				ctxt.Params = map[string]string{
					"registry":      "kind-registry.kind:5000",
					"builder-image": "localhost:5000/ods/buildah:latest",
					"tls-verify":    "false",
				}
			},
			WantRunSuccess: true,
			PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				checkResultingFiles(t, wsDir)
				checkResultingImage(t, ns, wsDir)
			},
		},
		"task should reuse existing image": {
			WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
			PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				tasktesting.InitAndCommitOrFatal(t, wsDir)
				tasktesting.WriteDotOdsOrFatal(t, wsDir, bitbucketProjectKey)
				buildAndPushImage(t, ns, wsDir)

				ctxt.Params = map[string]string{
					"registry":      "kind-registry.kind:5000",
					"builder-image": "localhost:5000/ods/buildah:latest",
					"tls-verify":    "false",
				}
			},
			WantRunSuccess: true,
			PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
				wsDir := ctxt.Workspaces["source"]
				checkResultingFiles(t, wsDir)
				checkResultingImage(t, ns, wsDir)
				checkLabelOnImage(t, ns, wsDir, "tasktestrun", "true")
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             "ClusterTask",          // could be read from task definition
				TaskName:                "ods-build-image-v0-1", // could be read from task definition
				Clients:                 c,
				Namespace:               ns,
				Timeout:                 5 * time.Minute, // depending on  the task we may need to increase or decrease it
				AlwaysKeepTmpWorkspaces: *alwaysKeepTmpWorkspacesFlag,
			})

		})

	}
}

// buildAndPushImage builds an image and pushes it to the registry.
// The used image tag equals the Git SHA that is being built, so the task
// will pick up the existing image.
// The image is labelled with "tasktestrun=true" so that it is possible to
// verify that the image has not been rebuild in the task.
func buildAndPushImage(t *testing.T, ns, wsDir string) {
	tag := getDockerImageTag(t, ns, wsDir)
	_, stderr, err := command.Run("docker", []string{
		"build", "--label", "tasktestrun=true", "-t", tag, filepath.Join(wsDir, "docker"),
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
		fmt.Sprintf(".ods/artifacts/image-digests/%s", filepath.Base(wsDir)),
	}
	for _, wf := range wantFiles {
		if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
			t.Fatalf("Want %s, but got nothing", wf)
		}
	}
}

func checkLabelOnImage(t *testing.T, ns, wsDir, wantLabelKey, wantLabelValue string) {
	stdout, stderr, err := command.Run("docker", []string{
		"image", "inspect", "--format", "{{ index .Config.Labels \"" + wantLabelKey + "\"}}",
		getDockerImageTag(t, ns, wsDir),
	})
	if err != nil {
		t.Fatalf("could not run get label on image: %s, stderr: %s", err, string(stderr))
	}
	got := strings.TrimSpace(string(stdout))
	if got != wantLabelValue {
		t.Fatalf("Want label %s=%s, but got value: %s", wantLabelKey, wantLabelValue, got)
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
	sha, err := getTrimmedFileContent(filepath.Join(wsDir, ".ods/git-commit-sha"))
	if err != nil {
		t.Fatalf("could not read git-commit-sha: %s", err)
	}
	return fmt.Sprintf("localhost:5000/%s/%s:%s", ns, filepath.Base(wsDir), sha)
}
