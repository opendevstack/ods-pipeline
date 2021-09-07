package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSPackageImage(t *testing.T) {
	runTaskTestCases(t,
		"ods-package-image",
		map[string]tasktesting.TestCase{
			"task should build image": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkResultingFiles(t, ctxt, wsDir)
					checkResultingImage(t, ctxt, wsDir)
				},
			},
			"task should reuse existing image": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					buildAndPushImageWithLabel(t, ctxt, wsDir)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkResultingFiles(t, ctxt, wsDir)
					checkResultingImage(t, ctxt, wsDir)
					checkLabelOnImage(t, ctxt, wsDir, "tasktestrun", "true")
				},
			},
		},
	)
}

// buildAndPushImageWithLabel builds an image and pushes it to the registry.
// The used image tag equals the Git SHA that is being built, so the task
// will pick up the existing image.
// The image is labelled with "tasktestrun=true" so that it is possible to
// verify that the image has not been rebuild in the task.
func buildAndPushImageWithLabel(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	tag := getDockerImageTag(t, ctxt, wsDir)
	t.Logf("Build image %s ahead of taskrun", tag)
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

func checkResultingFiles(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	wantFiles := []string{
		fmt.Sprintf(".ods/artifacts/image-digests/%s.json", ctxt.ODS.Component),
	}
	for _, wf := range wantFiles {
		if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
			t.Fatalf("Want %s, but got nothing", wf)
		}
	}
}

func checkLabelOnImage(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir, wantLabelKey, wantLabelValue string) {
	stdout, stderr, err := command.Run("docker", []string{
		"image", "inspect", "--format", "{{ index .Config.Labels \"" + wantLabelKey + "\"}}",
		getDockerImageTag(t, ctxt, wsDir),
	})
	if err != nil {
		t.Fatalf("could not run get label on image: %s, stderr: %s", err, string(stderr))
	}
	got := strings.TrimSpace(string(stdout))
	if got != wantLabelValue {
		t.Fatalf("Want label %s=%s, but got value: %s", wantLabelKey, wantLabelValue, got)
	}
}

func checkResultingImage(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	stdout, stderr, err := command.Run("docker", []string{
		"run", "--rm",
		getDockerImageTag(t, ctxt, wsDir),
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

func getDockerImageTag(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) string {
	sha, err := getTrimmedFileContent(filepath.Join(wsDir, ".ods/git-commit-sha"))
	if err != nil {
		t.Fatalf("could not read git-commit-sha: %s", err)
	}
	return fmt.Sprintf("localhost:5000/%s/%s:%s", ctxt.Namespace, ctxt.ODS.Component, sha)
}