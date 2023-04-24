package pipelinectxt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
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
	nexusURL := "https://nexus.example.com"
	permanentBaseURL := fmt.Sprintf("%s/%s%s/%s", nexusURL, nexus.TestPermanentRepository, group, artifactType)
	temporaryBaseURL := fmt.Sprintf("%s/%s%s/%s", nexusURL, nexus.TestPermanentRepository, group, artifactType)
	artifactsDir, err := os.MkdirTemp(".", "test-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(artifactsDir)
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}

	tests := map[string]struct {
		urls map[string][]string
		repo string
	}{
		"artifacts in permanent repo": {
			urls: map[string][]string{
				nexus.TestPermanentRepository: {
					permanentBaseURL + "/p1.txt", permanentBaseURL + "/p2.txt",
				},
			},
			repo: nexus.TestPermanentRepository,
		},
		"artifacts in temporary repo": {
			urls: map[string][]string{
				nexus.TestTemporaryRepository: {
					temporaryBaseURL + "/t1.txt", temporaryBaseURL + "/t2.txt",
				},
			},
			repo: nexus.TestTemporaryRepository,
		},
		"artifacts in no repo": {
			urls: map[string][]string{},
			repo: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			nexusClient.URLs = tc.urls
			am, err := DownloadGroup(nexusClient, tc.repo, group, artifactsDir, logger)
			if err != nil {
				t.Fatal(err)
			}
			if repoURLs, ok := tc.urls[tc.repo]; ok {
				for _, url := range repoURLs {
					ai := findArtifact(url, am.Artifacts)
					if ai == nil {
						t.Fatalf("expected artifact for %s in manifest, got none", url)
					}
					if ai.Directory != artifactType {
						t.Fatalf("want: %s, got: %s", artifactType, ai.Directory)
					}
					if !strings.HasSuffix(url, ai.Name) {
						t.Fatalf("want art to end in: %s, got: %s", ai.Name, url)
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

func findArtifact(url string, artifacts []ArtifactInfo) *ArtifactInfo {
	for _, a := range artifacts {
		if a.URL == url {
			return &a
		}
	}
	return nil
}
