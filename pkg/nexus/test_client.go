package nexus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	TestPermanentRepository = "ods-permanent-artifacts"
	TestTemporaryRepository = "ods-temporary-artifacts"
)

type TestArtifact struct {
	Path    string
	Content []byte
}

type TestClient struct {
	// Artifacts contains Artifacts per repository
	Artifacts map[string][]TestArtifact
}

// Download writes a dummy string into outfile.
func (c *TestClient) Download(url, outfile string) (int64, error) {
	for _, repo := range c.Artifacts {
		for _, ta := range repo {
			if ta.Path == url {
				return 0, os.WriteFile(outfile, ta.Content, 0644)
			}
		}
	}
	return 0, fmt.Errorf("url %s not found", url)
}

// Search responds with pre-registered URLs for repository.
// group is ignored.
func (c *TestClient) Search(repository, group string) ([]string, error) {
	urls := []string{}
	for _, ta := range c.Artifacts[repository] {
		if strings.HasPrefix(ta.Path, strings.TrimSuffix(group, "*")) {
			urls = append(urls, ta.Path)
		}
	}
	return urls, nil
}

// Upload is not needed for the test cases below.
func (c *TestClient) Upload(repository, group, file string) (string, error) {
	filename := filepath.Base(file)
	path := fmt.Sprintf("%s/%s", group, filename)
	c.Artifacts[repository] = append(c.Artifacts[repository], TestArtifact{Path: path})
	return fmt.Sprintf("%s/%s", repository, path), nil
}
