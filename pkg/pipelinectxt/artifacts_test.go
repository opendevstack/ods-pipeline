package pipelinectxt

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
)

func TestDownloadGroup(t *testing.T) {
	nexusClient := &nexus.TestClient{}
	repositories := []string{nexus.PermanentRepositoryDefault, nexus.TemporaryRepositoryDefault}
	group := "/my-project/my-repo/20d69cffd00080e20fa2f1419026a301cd0eecac"
	artifactType := "my-type"
	nexusURL := "https://nexus.example.com"
	permanentBaseURL := fmt.Sprintf("%s/%s%s/%s", nexusURL, nexus.PermanentRepositoryDefault, group, artifactType)
	temporaryBaseURL := fmt.Sprintf("%s/%s%s/%s", nexusURL, nexus.PermanentRepositoryDefault, group, artifactType)
	artifactsDir, err := ioutil.TempDir(".", "test-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(artifactsDir)
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}

	tests := map[string]struct {
		urls             map[string][]string
		wantSelectedRepo string
	}{
		"artifacts in permanent repo": {
			urls: map[string][]string{
				nexus.PermanentRepositoryDefault: {
					permanentBaseURL + "/p1.txt", permanentBaseURL + "/p2.txt",
				},
			},
			wantSelectedRepo: nexus.PermanentRepositoryDefault,
		},
		"artifacts in temporary repo": {
			urls: map[string][]string{
				nexus.TemporaryRepositoryDefault: {
					temporaryBaseURL + "/t1.txt", temporaryBaseURL + "/t2.txt",
				},
			},
			wantSelectedRepo: nexus.TemporaryRepositoryDefault,
		},
		"artifacts in both repos": {
			urls: map[string][]string{
				nexus.PermanentRepositoryDefault: {
					permanentBaseURL + "/p1.txt", permanentBaseURL + "/p2.txt",
				},
				nexus.TemporaryRepositoryDefault: {
					temporaryBaseURL + "/t1.txt", temporaryBaseURL + "/t2.txt",
				},
			},
			wantSelectedRepo: nexus.PermanentRepositoryDefault,
		},
		"artifacts in no repo": {
			urls:             map[string][]string{},
			wantSelectedRepo: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			nexusClient.URLs = tc.urls
			am, err := DownloadGroup(nexusClient, repositories, group, artifactsDir, logger)
			if err != nil {
				t.Fatal(err)
			}
			if am.SourceRepository != tc.wantSelectedRepo {
				t.Fatalf("want: %s, got: %s", tc.wantSelectedRepo, am.SourceRepository)
			}
			if repoURLs, ok := tc.urls[tc.wantSelectedRepo]; ok {
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
