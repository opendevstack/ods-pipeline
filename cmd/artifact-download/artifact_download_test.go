package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/gittest"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

func TestGetODSContextFromDir(t *testing.T) {
	sha := "7f96ec9fcf097e5b21687d402bc70370ac247d8a"
	dir, cleanup, err := gittest.CreateFakeGitRepoDir(
		"https://example.bitbucket.com/scm/ODS/ods-pipeline.git",
		"master",
		sha,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// As context is read from dir, no Bitbucket client should be required.
	ctxt, err := assembleODSContext("foo-cd", dir)
	if err != nil {
		t.Fatal(err)
	}
	wantContext := &pipelinectxt.ODSContext{
		Namespace:       "foo-cd",
		Project:         "ods",
		Repository:      "ods-pipeline",
		Component:       "pipeline",
		GitCommitSHA:    sha,
		GitFullRef:      "refs/heads/master",
		GitRef:          "master",
		GitURL:          "https://example.bitbucket.com/scm/ODS/ods-pipeline.git",
		PullRequestBase: "",
		PullRequestKey:  "",
	}
	if diff := cmp.Diff(wantContext, ctxt); diff != "" {
		t.Fatalf("context mismatch (-want +got):\n%s", diff)
	}
}

func TestRun(t *testing.T) {
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}
	project := "foo"
	repository := "bar"
	artifactType := "deployment"
	artifactName := "diff-dev.txt"
	repoCommitSHA := "7f96ec9fcf097e5b21687d402bc70370ac247d8a"
	subrepoName := "baz"
	subrepoCommitSHA := "8f96ec9fcf097e5b21687d402bc70370ac247d8a"

	// Git repo.
	dir, cleanup, err := gittest.CreateFakeGitRepoDir(
		fmt.Sprintf("https://example.bitbucket.com/scm/%s/%s.git", project, repository),
		"master", repoCommitSHA,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// Nexus client with corresponding artifact asset.
	nexusClient := &nexus.TestClient{
		Artifacts: map[string][]nexus.TestArtifact{
			nexus.TestPermanentRepository: {
				nexus.TestArtifact{
					Path: fmt.Sprintf(
						"%s/%s",
						nexus.ArtifactGroup(project, repository, repoCommitSHA, artifactType),
						artifactName,
					),
					Content: []byte("test"),
				},
				nexus.TestArtifact{
					Path: fmt.Sprintf(
						"%s/%s",
						nexus.ArtifactGroup(project, repository, repoCommitSHA, pipelinectxt.PipelineRunsDir),
						"foo-123.json",
					),
					Content: []byte(fmt.Sprintf("repositories: {%s: %s}", subrepoName, subrepoCommitSHA)),
				},
				nexus.TestArtifact{
					Path: fmt.Sprintf(
						"%s/%s",
						nexus.ArtifactGroup(project, subrepoName, subrepoCommitSHA, artifactType),
						artifactName,
					),
					Content: []byte("test"),
				},
			},
		},
	}

	// Temporary output directory.
	artifactsDir, err := os.MkdirTemp(".", "test-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(artifactsDir)

	// Program options.
	opts := options{
		namespace:       "foo-cd",
		outputDirectory: artifactsDir,
	}

	// Run main function and check for error / downloaded file.
	err = run(
		logger,
		opts,
		nexusClient,
		nexus.TestPermanentRepository,
		dir,
	)
	if err != nil {
		t.Fatal(err)
	}
	wantOutfile := filepath.Join(artifactsDir, repository, artifactType, artifactName)
	if _, err := os.Stat(wantOutfile); os.IsNotExist(err) {
		t.Fatalf("expected artifact downloaded to %s, got none", wantOutfile)
	}
	wantOutfile = filepath.Join(artifactsDir, subrepoName, artifactType, artifactName)
	if _, err := os.Stat(wantOutfile); os.IsNotExist(err) {
		t.Fatalf("expected artifact downloaded to %s, got none", wantOutfile)
	}
}
