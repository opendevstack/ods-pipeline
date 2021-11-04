package nexus

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ArtifactsManifest represents all downloaded artifacts.
type ArtifactsManifest struct {
	// SourceRepository identifies the repository artifacts where downloaded from
	SourceRepository string         `json:"sourceRepository"`
	Artifacts        []ArtifactInfo `json:"artifacts"`
}

// ArtifactInfo represents one artifact.
type ArtifactInfo struct {
	URL       string `json:"url"`
	Directory string `json:"directory"`
	Name      string `json:"name"`
}

// Download requests given url and writes the output into outfile.
func (c *Client) Download(url, outfile string) (int64, error) {
	if len(outfile) == 0 {
		outfile = path.Base(url)
	}
	out, err := os.Create(outfile)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	c.logger().Debugf("Download %s to %s", url, outfile)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Basic "+c.basicAuth())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return io.Copy(out, resp.Body)
}

// DownloadGroup searches given repositories in order for assets in given group.
// As soon as one repository has any asset in the group, the search is stopped
// and all fond artifacts are downloaded into artifactsDir.
// An artifacts manifest is returned describing the downloaded files.
// When none of the given repositories contains any artifacts under the group,
// no artifacts are downloaded and no error is returned.
func (c *Client) DownloadGroup(repositories []string, group, artifactsDir string) (*ArtifactsManifest, error) {
	// We want to target all artifacts underneath the group, hence the trailing '*'.
	nexusSearchGroup := fmt.Sprintf("%s/*", group)
	am := &ArtifactsManifest{
		Artifacts: []ArtifactInfo{},
	}
	sourceRepo, urls, err := c.searchForAssets(nexusSearchGroup, repositories)
	if err != nil {
		return nil, err
	}
	am.SourceRepository = sourceRepo

	for _, s := range urls {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		urlPathParts := strings.Split(u.Path, fmt.Sprintf("%s/", group))
		fileWithSubPath := urlPathParts[1]         // e.g. "pipeline-runs/foo-zh9gt0.json"
		aritfactName := path.Base(fileWithSubPath) // e.g. "pipeline-runs"
		artifactType := path.Dir(fileWithSubPath)  // e.g. "foo-zh9gt0.json"
		artifactsSubPath := filepath.Join(artifactsDir, artifactType)
		if _, err := os.Stat(artifactsSubPath); os.IsNotExist(err) {
			if err := os.MkdirAll(artifactsSubPath, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory: %s, error: %w", artifactsSubPath, err)
			}
		}
		outfile := filepath.Join(artifactsDir, fileWithSubPath)
		_, err = c.Download(s, outfile)
		if err != nil {
			return nil, err
		}
		am.Artifacts = append(am.Artifacts, ArtifactInfo{
			URL:       s,
			Directory: artifactType,
			Name:      aritfactName,
		})
	}
	return am, nil
}

// searchForAssets looks for assets in searchGroup for each repository in order.
// Once some assets are found, the repository and the found URLs are returned,
// skipping any further repositories that are given.
func (c *Client) searchForAssets(searchGroup string, repositories []string) (string, []string, error) {
	for _, r := range repositories {
		urls, err := c.Search(r, searchGroup)
		if err != nil {
			return "", nil, err
		}
		if len(urls) > 0 {
			c.logger().Infof("Found artifacts in repository %s inside group %s ...", r, searchGroup)
			return r, urls, nil
		}
		c.logger().Infof("No artifacts found in repository %s inside group %s.", r, searchGroup)
	}
	return "", []string{}, nil
}
