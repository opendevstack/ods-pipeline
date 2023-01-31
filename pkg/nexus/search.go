package nexus

import (
	"fmt"
	"strings"

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
		res = append(res, appendTestTLSPort(a.DownloadURL))
	}
	return res, nil
}

// appendTestTLSPort appends the TLS port used in testing to the URL.
func appendTestTLSPort(in string) string {
	return strings.Replace(
		strings.Replace(in, "localhost/", "localhost:8443/", -1),
		"ods-test-nexus-tls.kind/", "ods-test-nexus-tls.kind:8443/", -1)
}
