package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/internal/gittest"
	"github.com/opendevstack/pipeline/internal/installation"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
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

	opts := options{
		namespace: "foo-cd",
		tag:       pipelinectxt.WIP, // for WIP version, ODS context is from dir
	}
	// As context is read from dir, no Bitbucket client should be required.
	ctxt, err := getODSContext(opts, nil, dir)
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
		Version:         "WIP",
		Environment:     "",
		PullRequestBase: "",
		PullRequestKey:  "",
	}
	if diff := cmp.Diff(wantContext, ctxt); diff != "" {
		t.Fatalf("context mismatch (-want +got):\n%s", diff)
	}
}

func TestGetODSContextFromBitbucketRepo(t *testing.T) {
	sha := "7f96ec9fcf097e5b21687d402bc70370ac247d8a"
	opts := options{
		namespace:  "foo-cd",
		project:    "foo",
		repository: "bar",
		tag:        "v1.0.0",
	}

	bitbucketClient := &bitbucket.TestClient{
		Tags: []bitbucket.Tag{
			{
				DisplayID:    "v1.0.0",
				LatestCommit: sha,
			},
		},
	}
	// As context is read from fake Bitbucket, dir value is unused.
	ctxt, err := getODSContext(opts, bitbucketClient, ".")
	if err != nil {
		t.Fatal(err)
	}
	wantContext := &pipelinectxt.ODSContext{
		Namespace:    "foo-cd",
		Project:      "foo",
		Repository:   "bar",
		GitCommitSHA: sha,
	}
	if diff := cmp.Diff(wantContext, ctxt); diff != "" {
		t.Fatalf("context mismatch (-want +got):\n%s", diff)
	}
}

func TestGetSubrepoODSContext(t *testing.T) {
	ctxt := &pipelinectxt.ODSContext{
		Namespace:    "foo-cd",
		Project:      "foo",
		Repository:   "bar",
		GitCommitSHA: "7f96ec9fcf097e5b21687d402bc70370ac247d8a",
	}

	tests := map[string]struct {
		opts     options
		subrepo  config.Repository
		branches []bitbucket.Branch
		tags     []bitbucket.Tag
		wantCtxt *pipelinectxt.ODSContext
	}{
		"tag given": {
			opts: options{
				namespace:  "foo-cd",
				project:    "foo",
				repository: "bar",
				tag:        "v1.0.0",
			},
			subrepo: config.Repository{Name: "baz"},
			tags: []bitbucket.Tag{
				{
					DisplayID:    "v1.0.0",
					LatestCommit: "f31532481bf29dffe02367f050f4a3f4dd7845ed",
				},
			},
			wantCtxt: &pipelinectxt.ODSContext{
				Namespace:    "foo-cd",
				Project:      "foo",
				Repository:   "baz",
				GitCommitSHA: "f31532481bf29dffe02367f050f4a3f4dd7845ed",
			},
		},
		"WIP given and no configured branch": {
			opts: options{
				namespace:  "foo-cd",
				project:    "foo",
				repository: "bar",
				tag:        pipelinectxt.WIP,
			},
			subrepo: config.Repository{Name: "baz"},
			branches: []bitbucket.Branch{
				{
					ID:           "refs/heads/master",
					LatestCommit: "af31532481bf29dffe02367f050f4a3f4dd7845ed",
				},
			},
			wantCtxt: &pipelinectxt.ODSContext{
				Namespace:    "foo-cd",
				Project:      "foo",
				Repository:   "baz",
				GitCommitSHA: "af31532481bf29dffe02367f050f4a3f4dd7845ed",
			},
		},
		"WIP given and configured branch": {
			opts: options{
				namespace:  "foo-cd",
				project:    "foo",
				repository: "bar",
				tag:        pipelinectxt.WIP,
			},
			subrepo: config.Repository{Name: "baz", Branch: "production"},
			branches: []bitbucket.Branch{
				{
					ID:           "refs/heads/master",
					LatestCommit: "af31532481bf29dffe02367f050f4a3f4dd7845ed",
				},
				{
					ID:           "refs/heads/production",
					LatestCommit: "bf31532481bf29dffe02367f050f4a3f4dd7845ed",
				},
			},
			wantCtxt: &pipelinectxt.ODSContext{
				Namespace:    "foo-cd",
				Project:      "foo",
				Repository:   "baz",
				GitCommitSHA: "bf31532481bf29dffe02367f050f4a3f4dd7845ed",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bitbucketClient := &bitbucket.TestClient{
				Branches: tc.branches,
				Tags:     tc.tags,
			}
			got, err := getSubrepoODSContext(ctxt, tc.subrepo, tc.opts, bitbucketClient)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.wantCtxt, got); diff != "" {
				t.Fatalf("context mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRun(t *testing.T) {
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}
	project := "foo"
	repository := "bar"
	tag := "v1.0.0"
	commitSHA := "f31532481bf29dffe02367f050f4a3f4dd7845ed"
	artifactType := "deployment"
	artifactName := "diff-dev.txt"

	// Bitbucket client with corresponding Git tag and empty ods.yaml.
	bitbucketClient := &bitbucket.TestClient{
		Tags: []bitbucket.Tag{
			{
				DisplayID:    tag,
				LatestCommit: commitSHA,
			},
		},
		Files: map[string][]byte{
			"ods.yaml": []byte("repositories: []"),
		},
	}

	// Nexus client with corresponding artifact asset.
	nexusClient := &nexus.TestClient{
		URLs: map[string][]string{
			nexus.PermanentRepositoryDefault: {
				fmt.Sprintf(
					"https://nexus.example.com/%s%s/%s",
					nexus.PermanentRepositoryDefault,
					nexus.ArtifactGroup(project, repository, commitSHA, artifactType),
					artifactName,
				),
			},
		},
	}

	// Temporary output directory.
	artifactsDir, err := ioutil.TempDir(".", "test-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(artifactsDir)

	// Program options.
	opts := options{
		namespace:       "foo-cd",
		project:         project,
		repository:      repository,
		tag:             tag,
		outputDirectory: artifactsDir,
	}

	// Run main function and check for error / downloaded file.
	err = run(
		logger,
		opts,
		nexusClient,
		&installation.NexusRepositories{
			Permanent: nexus.PermanentRepositoryDefault,
			Temporary: nexus.TemporaryRepositoryDefault,
		},
		bitbucketClient,
		".",
	)
	if err != nil {
		t.Fatal(err)
	}
	wantOutfile := filepath.Join(artifactsDir, tag, repository, artifactType, artifactName)
	if _, err := os.Stat(wantOutfile); os.IsNotExist(err) {
		t.Fatalf("expected artifact downloaded to %s, got none", wantOutfile)
	}
}
