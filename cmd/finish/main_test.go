package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

func TestUploadArtifacts(t *testing.T) {
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}
	nexusRepo := "foo"
	nexusClient := &nexus.TestClient{
		URLs: map[string][]string{},
	}
	tempWorkingDir, cleanup, err := prepareTempWorkingDir(nexusRepo)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	ctxt := &pipelinectxt.ODSContext{
		Version:      "1.0.0",
		Project:      "my-project",
		Repository:   "my-repo",
		GitCommitSHA: "8d351a10fb428c0c1239530256e21cf24f136e73",
	}
	t.Log("Write dummy artifact")
	artifactsDir := filepath.Join(tempWorkingDir, pipelinectxt.ArtifactsPath)
	err = writeArtifactFile(artifactsDir, pipelinectxt.ImageDigestsDir, "foo.json")
	if err != nil {
		t.Fatal(err)
	}

	err = uploadArtifacts(logger, nexusClient, nexusRepo, tempWorkingDir, ctxt)
	if err != nil {
		t.Fatal(err)
	}
	if len(nexusClient.URLs[nexusRepo]) != 1 {
		t.Fatalf("want 1 uploaded file, got: %v", nexusClient.URLs[nexusRepo])
	}
	wantFile := "/my-project/my-repo/8d351a10fb428c0c1239530256e21cf24f136e73/image-digests/foo.json"
	if !nexusRepoContains(nexusClient.URLs[nexusRepo], wantFile) {
		t.Fatalf("want: %s, got: %s", wantFile, nexusClient.URLs[nexusRepo][0])
	}
}

func TestHandleArtifacts(t *testing.T) {
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}
	nexusRepo := "temporary"
	nexusClient := &nexus.TestClient{
		URLs: map[string][]string{},
	}
	tempWorkingDir, cleanup, err := prepareTempWorkingDir(nexusRepo)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	ctxt := &pipelinectxt.ODSContext{
		Version:      "1.0.0",
		Project:      "my-project",
		Repository:   "my-repo",
		GitCommitSHA: "8d351a10fb428c0c1239530256e21cf24f136e73",
	}
	opts := options{
		aggregateTasksStatus:     "Succeeded",
		nexusTemporaryRepository: "temporary",
		nexusPermanentRepository: "permanent",
		pipelineRunName:          "foo",
	}
	t.Log("Add ods.yaml")
	subrepoName := "my-subrepo"
	ods := config.ODS{Repositories: []config.Repository{
		{
			Name: subrepoName,
		},
	}}
	out, err := json.Marshal(ods)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(tempWorkingDir, "ods.yaml"), out, 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Add subrepository directory")
	subrepoDir := filepath.Join(tempWorkingDir, pipelinectxt.SubreposPath, subrepoName)
	err = os.MkdirAll(subrepoDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Write pipeline context into subrepository directory")
	err = writeTestContext("foo-cd", subrepoDir, subrepoName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Write empty artifacts manifest for subrepository")
	subrepoArtifactsDir := filepath.Join(subrepoDir, pipelinectxt.ArtifactsPath)
	err = writeArtifactsManifest(subrepoArtifactsDir, nexusRepo)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Write dummy artifact for subrepository")
	err = writeArtifactFile(subrepoArtifactsDir, pipelinectxt.ImageDigestsDir, "bar.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = handleArtifacts(logger, nexusClient, opts, tempWorkingDir, ctxt)
	if err != nil {
		t.Fatal(err)
	}
	if len(nexusClient.URLs[nexusRepo]) != 2 {
		t.Fatalf("want 2 uploaded files, got: %v", nexusClient.URLs[nexusRepo])
	}
	wantFile := "/my-project/my-repo/8d351a10fb428c0c1239530256e21cf24f136e73/pipeline-runs/foo.json"
	if !nexusRepoContains(nexusClient.URLs[nexusRepo], wantFile) {
		t.Fatalf("want: %s, got: %s", wantFile, nexusClient.URLs[nexusRepo])
	}
	wantFileSubrepo := "/my-project/my-subrepo/86mz0pa4ci0ke5o27gmdnlnqdwdvkrx1iw8mdpta/image-digests/bar.txt"
	if !nexusRepoContains(nexusClient.URLs[nexusRepo], wantFileSubrepo) {
		t.Fatalf("want: %s, got: %s", wantFileSubrepo, nexusClient.URLs[nexusRepo])
	}
}

// nexusRepoContains checks if haystack contains needle
func nexusRepoContains(haystack []string, needle string) bool {
	for _, f := range haystack {
		if f == needle {
			return true
		}
	}
	return false
}

// prepareTempWorkingDir creates a temporary directory which includes an
// artifacts manifest file. The returned function should be used for cleanup.
func prepareTempWorkingDir(nexusRepo string) (string, func(), error) {
	cleanup := func() {}
	tempWorkingDir, err := ioutil.TempDir(".", "test-upload-")
	if err != nil {
		return tempWorkingDir, cleanup, err
	}
	cleanup = func() { os.RemoveAll(tempWorkingDir) }
	artifactsDir := filepath.Join(tempWorkingDir, pipelinectxt.ArtifactsPath)
	err = writeArtifactsManifest(artifactsDir, nexusRepo)
	if err != nil {
		return tempWorkingDir, cleanup, err
	}

	return tempWorkingDir, cleanup, err
}

// writeTestContext writes an ODS context into wsDir.
func writeTestContext(ns, wsDir, repoName string) error {
	ctxt := &pipelinectxt.ODSContext{
		Namespace:    ns,
		Project:      "my-project",
		Repository:   repoName,
		Component:    repoName,
		GitCommitSHA: "86mz0pa4ci0ke5o27gmdnlnqdwdvkrx1iw8mdpta",
		GitFullRef:   "refs/heads/master",
		GitRef:       "master",
		GitURL:       "http://bitbucket.acme.org/scm/my-project/my-repo.git",
		Environment:  "dev",
		Version:      pipelinectxt.WIP,
	}
	return ctxt.WriteCache(wsDir)
}

// writeArtifactFile writes a dummy file named filename into artifactsDir/subdir.
func writeArtifactFile(artifactsDir, subdir, filename string) error {
	artifactsSubDir := filepath.Join(artifactsDir, subdir)
	err := os.MkdirAll(artifactsSubDir, 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(artifactsSubDir, filename), []byte("test"), 0644)
}

// writeArtifactsManifest writes an artigact manifest JSON file into artifactsDir.
func writeArtifactsManifest(artifactsDir, nexusRepo string) error {
	am := &pipelinectxt.ArtifactsManifest{
		SourceRepository: nexusRepo,
		Artifacts:        []pipelinectxt.ArtifactInfo{},
	}
	return pipelinectxt.WriteJsonArtifact(am, artifactsDir, pipelinectxt.ArtifactsManifestFilename)
}
