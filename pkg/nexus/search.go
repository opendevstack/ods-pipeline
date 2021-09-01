package nexus

import (
	"fmt"

	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

// Search gets URLs of assets in given repository group.
func (c *Client) Search(repository, group string) ([]string, error) {
	c.logger().Debugf("Search for %s in %s", group, repository)
	query := nexusrm.NewSearchQueryBuilder().Repository(repository).Group(group)
	assets, err := nexusrm.SearchAssets(c.rm, query)
	if err != nil {
		return nil, fmt.Errorf("could not search assets: %w", err)
	}

	res := []string{}
	for _, a := range assets {
		res = append(res, a.DownloadURL)
	}
	return res, nil
}
