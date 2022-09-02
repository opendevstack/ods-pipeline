package nexus

import (
	"fmt"
	"os"
	"path/filepath"
)

type TestClient struct {
	// URLs contains URLs per repository
	URLs map[string][]string
}

// Download writes a dummy string into outfile.
func (c *TestClient) Download(url, outfile string) (int64, error) {
	return 0, os.WriteFile(outfile, []byte("test"), 0644)
}

// Search responds with pre-registered URLs for repository.
// group is ignored.
func (c *TestClient) Search(repository, group string) ([]string, error) {
	return c.URLs[repository], nil
}

// Upload is not needed for the test cases below.
func (c *TestClient) Upload(repository, group, file string) (string, error) {
	filename := filepath.Base(file)
	path := fmt.Sprintf("%s/%s", group, filename)
	c.URLs[repository] = append(
		c.URLs[repository],
		path,
	)
	return fmt.Sprintf("%s/%s", repository, path), nil
}
