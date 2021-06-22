package nexus

import (
	"fmt"

	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

// Search gets URLs
func Search(URL, user, password, repository, group string) ([]string, error) {
	rm, err := nexusrm.New(
		URL,
		user,
		password,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create nexus client: %w", err)
	}

	query := nexusrm.NewSearchQueryBuilder().Repository(repository).Group(group)
	assets, err := nexusrm.SearchAssets(rm, query)
	if err != nil {
		return nil, fmt.Errorf("could not search assets: %w", err)
	}

	res := []string{}
	for _, a := range assets {
		res = append(res, a.DownloadURL)
	}
	return res, nil
}
