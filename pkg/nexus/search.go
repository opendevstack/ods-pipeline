package nexus

import (
	"fmt"

	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

// Search gets URLs of assets in given repository group.
func (c *Client) Search(group string) ([]string, error) {
	query := nexusrm.NewSearchQueryBuilder().Repository(c.Repository).Group(group)
	assets, err := nexusrm.SearchAssets(c.RM, query)
	if err != nil {
		return nil, fmt.Errorf("could not search assets: %w", err)
	}

	res := []string{}
	for _, a := range assets {
		res = append(res, a.DownloadURL)
	}
	return res, nil
}
