package nexus

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type TestClient struct {
	// URLs contains URLs per repository
	URLs map[string][]string
}

// Download writes a dummy string into outfile.
func (c *TestClient) Download(url, outfile string) (int64, error) {
	return 0, ioutil.WriteFile(outfile, []byte("test"), 0644)
}

// Search responds with pre-registered URLs for repository.
// group is ignored.
func (c *TestClient) Search(repository, group string) ([]string, error) {
	return c.URLs[repository], nil
}

// Upload is not needed for the test cases below.
func (c *TestClient) Upload(repository, group, file string) error {
	filename := filepath.Base(file)
	c.URLs[repository] = append(
		c.URLs[repository],
		fmt.Sprintf("%s/%s", group, filename),
	)
	return nil
}
