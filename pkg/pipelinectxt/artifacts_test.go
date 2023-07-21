package pipelinectxt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/ods-pipeline/pkg/logging"
	"github.com/opendevstack/ods-pipeline/pkg/nexus"
)

func TestReadArtifactsDir(t *testing.T) {
	artifactsDir, err := os.MkdirTemp(".", "test-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(artifactsDir)
	artifactsSubDir := filepath.Join(artifactsDir, PipelineRunsDir)
	err = os.MkdirAll(artifactsSubDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	artifactsFile := filepath.Join(artifactsSubDir, "foo.txt")
	err = os.WriteFile(artifactsFile, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	m, err := ReadArtifactsDir(artifactsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 1 {
		t.Fatalf("want 1 entry, got: %v", m)
	}
	if m[PipelineRunsDir][0] != "foo.txt" {
		t.Fatalf("want foo.txt, got: %v", m)
	}
}

func TestDownloadGroup(t *testing.T) {
	nexusClient := &nexus.TestClient{}
	group := "/my-project/my-repo/20d69cffd00080e20fa2f1419026a301cd0eecac"
	artifactType := "my-type"
	basePath := fmt.Sprintf("%s/%s", group, artifactType)
	artifactsDir, err := os.MkdirTemp(".", "test-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(artifactsDir)
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}

	tests := map[string]struct {
		urls map[string][]nexus.TestArtifact
		repo string
	}{
		"artifacts in permanent repo": {
			urls: map[string][]nexus.TestArtifact{
				nexus.TestPermanentRepository: {
					nexus.TestArtifact{
						Path:    basePath + "/p1.txt",
						Content: []byte("test"),
					}, nexus.TestArtifact{
						Path:    basePath + "/p2.txt",
						Content: []byte("test"),
					},
				},
			},
			repo: nexus.TestPermanentRepository,
		},
		"artifacts in temporary repo": {
			urls: map[string][]nexus.TestArtifact{
				nexus.TestTemporaryRepository: {
					nexus.TestArtifact{
						Path:    basePath + "/t1.txt",
						Content: []byte("test"),
					}, nexus.TestArtifact{
						Path:    basePath + "/t2.txt",
						Content: []byte("test"),
					},
				},
			},
			repo: nexus.TestTemporaryRepository,
		},
		"artifacts in no repo": {
			urls: map[string][]nexus.TestArtifact{},
			repo: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			nexusClient.Artifacts = tc.urls
			am, err := DownloadGroup(nexusClient, tc.repo, group, artifactsDir, logger)
			if err != nil {
				t.Fatal(err)
			}
			if repoURLs, ok := tc.urls[tc.repo]; ok {
				for _, ta := range repoURLs {
					ai := findArtifact(ta.Path, am.Artifacts)
					if ai == nil {
						t.Fatalf("expected artifact for %s in manifest, got none", ta)
					}
					if ai.Directory != artifactType {
						t.Fatalf("want: %s, got: %s", artifactType, ai.Directory)
					}
					if !strings.HasSuffix(ta.Path, ai.Name) {
						t.Fatalf("want art to end in: %s, got: %s", ai.Name, ta)
					}
					wantOutfile := filepath.Join(artifactsDir, ai.Directory, ai.Name)
					if _, err := os.Stat(wantOutfile); os.IsNotExist(err) {
						t.Fatalf("expected artifact downloaded to %s, got none", wantOutfile)
					}
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := map[string]struct {
		manifest *ArtifactsManifest
		repo     string
		dir      string
		name     string
		want     bool
	}{
		"different repo": {
			manifest: &ArtifactsManifest{
				Repository: "a",
				Artifacts: []ArtifactInfo{
					{
						Directory: "b",
						Name:      "c",
					},
				},
			},
			repo: "x",
			dir:  "b",
			name: "c",
			want: false,
		},
		"same repo, different dir": {
			manifest: &ArtifactsManifest{
				Repository: "a",
				Artifacts: []ArtifactInfo{
					{
						Directory: "b",
						Name:      "c",
					},
				},
			},
			repo: "a",
			dir:  "x",
			name: "c",
			want: false,
		},
		"same repo, same dir, different name": {
			manifest: &ArtifactsManifest{
				Repository: "a",
				Artifacts: []ArtifactInfo{
					{
						Directory: "b",
						Name:      "c",
					},
				},
			},
			repo: "a",
			dir:  "b",
			name: "x",
			want: false,
		},
		"match": {
			manifest: &ArtifactsManifest{
				Repository: "a",
				Artifacts: []ArtifactInfo{
					{
						Directory: "b",
						Name:      "c",
					},
				},
			},
			repo: "a",
			dir:  "b",
			name: "c",
			want: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.manifest.Contains(tc.repo, tc.dir, tc.name) != tc.want {
				t.Errorf("Want %q to contain=%v, but did not", name, tc.want)
			}
		})
	}
}

func findArtifact(url string, artifacts []ArtifactInfo) *ArtifactInfo {
	for _, a := range artifacts {
		if a.URL == url {
			return &a
		}
	}
	return nil
}
