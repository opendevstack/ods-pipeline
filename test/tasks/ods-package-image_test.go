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
	"golang.org/x/exp/slices"
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
			"task should build image with additional tags": {
				WorkspaceDirMapping: map[string]string{"source": "hello-world-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
					ctxt.Params = map[string]string{
						"extra-tags": "'latest cool'",
					}
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]
					checkResultingFiles(t, ctxt, wsDir)
					checkTagFiles(t, ctxt, wsDir, []string{"latest", "cool"})
					checkTags(t, ctxt, wsDir, []string{ctxt.ODS.GitCommitSHA, "latest", "cool"})
					checkResultingImageHelloWorld(t, ctxt, wsDir)
					checkTaggedImageHelloWorld(t, ctxt, wsDir, "latest")
					checkTaggedImageHelloWorld(t, ctxt, wsDir, "cool")
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
		fmt.Sprintf(".ods/artifacts/sboms/%s.spdx", ctxt.ODS.Component),
	}
	for _, wf := range wantFiles {
		if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
			t.Fatalf("Want %s, but got nothing", wf)
		}
	}
}

func checkTagFiles(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string, tags []string) {
	wantFiles := []string{}
	for _, tag := range tags {
		wantFiles = append(wantFiles, fmt.Sprintf(".ods/artifacts/image-digests/%s-%s.json", ctxt.ODS.Component, tag))
	}
	for _, wf := range wantFiles {
		if _, err := os.Stat(filepath.Join(wsDir, wf)); os.IsNotExist(err) {
			t.Fatalf("Want %s, but got nothing", wf)
		}
	}
}

func checkTags(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string, expectedTags []string) {
	// registry := "kind-registry.kind:5000"
	registry := "localhost:5000"
	args := []string{
		"inspect",
		`--format={{.RepoTags}}`,
	}
	imageNsStreamSha := fmt.Sprintf("%s/%s:%s", ctxt.Namespace, ctxt.ODS.Component, ctxt.ODS.GitCommitSHA)
	imageRef := fmt.Sprintf("%s/%s", registry, imageNsStreamSha)
	args = append(args, imageRef)

	stdout, _, err := command.RunBuffered("docker", args)
	if err != nil {
		t.Fatalf("docker %s: %s", fmt.Sprint(args), err)
	}
	tags, err := parseDockerInspectDigestTags(string(stdout))
	if err != nil {
		t.Fatalf("parse tags failed: %s", err)
	}
	for _, expectedTag := range expectedTags {
		if !slices.Contains(tags, expectedTag) {
			t.Fatalf("Expected tags=%s to be in actual tags=%s", fmt.Sprint(expectedTags), fmt.Sprint(tags))
		}
	}
}

func parseDockerInspectDigestTags(out string) ([]string, error) {
	t := strings.TrimSpace(out)
	if !(strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]")) {
		return nil, fmt.Errorf("docker inspect: unexpected tag response expecting tags to be in brackets %s", t)
	}
	t = t[1 : len(t)-1]
	// expecting t to have space separated tags.
	tags := strings.Split(t, " ")
	return tags, nil
}

func runSpecifiedImage(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string, image string) string {
	stdout, stderr, err := command.RunBuffered("docker", []string{
		"run", "--rm",
		image,
	})
	if err != nil {
		t.Fatalf("could not run built image: %s, stderr: %s", err, string(stderr))
	}
	got := strings.TrimSpace(string(stdout))
	return got
}

func runResultingImage(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) string {
	got := runSpecifiedImage(t, ctxt, wsDir, getDockerImageTag(t, ctxt, wsDir))
	return got
}

func checkResultingImageHelloWorld(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	got := runResultingImage(t, ctxt, wsDir)
	want := "Hello World"
	if got != want {
		t.Fatalf("Want %s, but got %s", want, got)
	}
}

func checkTaggedImageHelloWorld(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string, tag string) {
	image := fmt.Sprintf("localhost:5000/%s/%s:%s", ctxt.Namespace, ctxt.ODS.Component, tag)
	got := runSpecifiedImage(t, ctxt, wsDir, image)
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
	generateImageArtifact(t, ctxt, tag, wsDir)
	generateImageSBOMArtifact(t, ctxt, wsDir)
}

func generateImageArtifact(t *testing.T, ctxt *tasktesting.TaskRunContext, tag string, wsDir string) {
	t.Logf("Generating image artifact")
	sha, err := getTrimmedFileContent(filepath.Join(wsDir, ".ods/git-commit-sha"))
	if err != nil {
		t.Fatalf("could not read git-commit-sha: %s", err)
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
	if err != nil {
		t.Fatalf("could not create image artifact: %s", err)
	}
}

func generateImageSBOMArtifact(t *testing.T, ctxt *tasktesting.TaskRunContext, wsDir string) {
	t.Logf("Generating image SBOM artifact")
	artifactsDir := filepath.Join(wsDir, pipelinectxt.SBOMsPath)
	sbomArtifactFilename := fmt.Sprintf("%s.%s", ctxt.ODS.Component, pipelinectxt.SBOMsFormat)
	err := os.MkdirAll(artifactsDir, 0755)
	if err != nil {
		t.Fatalf("could not create %s: %s", artifactsDir, err)
	}
	_, err = os.Create(filepath.Join(artifactsDir, sbomArtifactFilename))
	if err != nil {
		t.Fatalf("could not create image SBOM artifact: %s", err)
	}
}
