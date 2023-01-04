package tasks

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/installation"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
)

func TestTaskODSPackageImage(t *testing.T) {
	runTaskTestCases(t,
		"ods-package-image",
		[]tasktesting.Service{
			tasktesting.Nexus,
		},
		map[string]tasktesting.TestCase{
			"task should build image and use nexus args": {
				WorkspaceDirMapping: map[string]string{"source": "hello-nexus-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkResultingFiles(t, ctxt, wsDir)
					checkResultingImageHelloNexus(t, ctxt, wsDir)
				},
			},
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
					checkResultingImageHelloWorld(t, ctxt, wsDir)
				},
			},
			"task should reuse existing image": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					tag := getDockerImageTag(t, ctxt, wsDir)
					generateArtifacts(t, ctxt, tag, wsDir)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkResultingFiles(t, ctxt, wsDir)
				},
			},
			"task should build image with build extra args param": {
				WorkspaceDirMapping: map[string]string{"source": "hello-build-extra-args-app"},
				TaskParamsMapping:   map[string]string{"buildah-build-extra-args": "'--build-arg=firstArg=one --build-arg=secondArg=two'"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkResultingFiles(t, ctxt, wsDir)
					checkResultingImageHelloBuildExtraArgs(t, ctxt, wsDir)
				},
			},
		},
	)
}

func checkResultingFiles(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	wantFiles := []string{
		fmt.Sprintf(".ods/artifacts/image-digests/%s.json", ctxt.ODS.Component),
		".ods/artifacts/sboms/sbom.spdx",
	}
	for _, wf := range wantFiles {
		if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
			t.Fatalf("Want %s, but got nothing", wf)
		}
	}
}

func runResultingImage(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) string {
	stdout, stderr, err := command.RunBuffered("docker", []string{
		"run", "--rm",
		getDockerImageTag(t, ctxt, wsDir),
	})
	if err != nil {
		t.Fatalf("could not run built image: %s, stderr: %s", err, string(stderr))
	}
	got := strings.TrimSpace(string(stdout))
	return got
}

func checkResultingImageHelloWorld(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	got := runResultingImage(t, ctxt, wsDir)
	want := "Hello World"
	if got != want {
		t.Fatalf("Want %s, but got %s", want, got)
	}
}

func checkResultingImageHelloNexus(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	got := runResultingImage(t, ctxt, wsDir)
	gotLines := strings.Split(got, "\n")

	ncc, err := installation.NewNexusClientConfig(
		ctxt.Clients.KubernetesClientSet, ctxt.Namespace, &logging.LeveledLogger{Level: logging.LevelDebug},
	)
	if err != nil {
		t.Fatalf("could not create Nexus client config: %s", err)
	}

	// nexusClient := tasktesting.NexusClientOrFatal(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)
	nexusUrlString := string(ncc.BaseURL)
	nexusUrl, err := url.Parse(nexusUrlString)
	if err != nil {
		t.Fatalf("could not determine nexusUrl from nexusClient: %s", err)
	}

	wantUsername := "developer"
	if ncc.Username != wantUsername {
		t.Fatalf("Want %s, but got %s", wantUsername, ncc.Username)
	}

	wantSecret := "s3cr3t"
	if ncc.Password != wantSecret {
		t.Fatalf("Want %s, but got %s", wantSecret, ncc.Password)
	}

	want := []string{
		fmt.Sprintf("nexusUrl=%s", nexusUrlString),
		fmt.Sprintf("nexusUsername=%s", ncc.Username),
		fmt.Sprintf("nexusPassword=%s", ncc.Password),
		fmt.Sprintf("nexusAuth=%s:%s", ncc.Username, ncc.Password),
		fmt.Sprintf("nexusUrlWithAuth=http://%s:%s@%s", ncc.Username, ncc.Password, nexusUrl.Host),
		fmt.Sprintf("nexusHost=%s", nexusUrl.Host),
	}
	if diff := cmp.Diff(want, gotLines); diff != "" {
		t.Fatalf("context mismatch (-want +got):\n%s", diff)
	}
}

func checkResultingImageHelloBuildExtraArgs(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	got := runResultingImage(t, ctxt, wsDir)
	gotLines := strings.Split(got, "\n")

	want := []string{
		fmt.Sprintf("firstArg=%s", "one"),
		fmt.Sprintf("secondArg=%s", "two"),
	}
	if diff := cmp.Diff(want, gotLines); diff != "" {
		t.Fatalf("context mismatch (-want +got):\n%s", diff)
	}
}

func getDockerImageTag(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) string {
	sha, err := getTrimmedFileContent(filepath.Join(wsDir, ".ods/git-commit-sha"))
	if err != nil {
		t.Fatalf("could not read git-commit-sha: %s", err)
	}
	return fmt.Sprintf("localhost:5000/%s/%s:%s", ctxt.Namespace, ctxt.ODS.Component, sha)
}

func generateArtifacts(t *testing.T, ctxt *tasktesting.TaskRunContext, tag string, wsDir string) {
	t.Logf("Generating artifacts for image %s", tag)
	t.Logf("Generating image artifact")
	err := generateImageArtifact(ctxt, tag, wsDir)
	if err != nil {
		t.Fatalf("could not create image artifact: %s", err)
	}
	t.Logf("Generating image SBOM artifact")
	err = generateImageSBOMArtifact(ctxt, wsDir)
	if err != nil {
		t.Fatalf("could not create image SBOM artifact: %s", err)
	}
}

func generateImageArtifact(ctxt *tasktesting.TaskRunContext, tag string, wsDir string) error {
	sha, err := getTrimmedFileContent(filepath.Join(wsDir, ".ods/git-commit-sha"))
	if err != nil {
		return fmt.Errorf("could not read git-commit-sha: %s", err)
	}
	ia := artifact.Image{
		Ref:        tag,
		Registry:   "kind-registry.kind:5000",
		Repository: ctxt.Namespace,
		Name:       ctxt.ODS.Component,
		Tag:        sha,
		Digest:     "abc",
	}
	imageArtifactFilename := fmt.Sprintf("%s.json", ctxt.ODS.Component)
	err = pipelinectxt.WriteJsonArtifact(ia, filepath.Join(wsDir, pipelinectxt.ImageDigestsPath), imageArtifactFilename)
	return err
}

func generateImageSBOMArtifact(ctxt *tasktesting.TaskRunContext, wsDir string) error {
	artifactsDir := filepath.Join(wsDir, pipelinectxt.SbomsPath)
	// imageArtifactFilename := fmt.Sprintf("%s.json", p.ctxt.Component)
	sbomArtifactFilename := "sbom.spdx"
	err := os.MkdirAll(artifactsDir, 0755)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", artifactsDir, err)
	}
	_, err = os.Create(filepath.Join(artifactsDir, sbomArtifactFilename))
	if err != nil {
		return fmt.Errorf("could not create SBOM fake file: %s", err)
	}
	return err
}
