package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Branch struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	IsDefault       bool   `json:"isDefault"`
}

type BranchPage struct {
	Size       int      `json:"size"`
	Limit      int      `json:"limit"`
	IsLastPage bool     `json:"isLastPage"`
	Values     []Branch `json:"values"`
	Start      int      `json:"start"`
}

type BranchListParams struct {
	// Base is the base branch or tag to compare each branch to (for the
	// metadata providers that uses that information).
	Base string `json:"base"`
	// Details controls whether to retrieve plugin-provided metadata about each branch.
	Details bool `json:"details"`
	// FilterText is the the text to match on. The match seems to be a prefix match.
	FilterText string `json:"filterText"`
	// OrderBy determines ordering of refs.
	// Either ALPHABETICAL (by name) or MODIFICATION (last updated).
	OrderBy string `json:"orderBy"`
	// BoostMatches controls whether exact and prefix matches will be boosted to the top
	BoostMatches bool `json:"boostMatches"`
}

type BranchClientInterface interface {
	BranchList(projectKey string, repositorySlug string, params BranchListParams) (*BranchPage, error)
}

// BranchList retrieves the branches matching the supplied filterText param.
// The authenticated user must have REPO_READ permission for the specified repository to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.14.0/bitbucket-rest.html#idp211
func (c *Client) BranchList(projectKey string, repositorySlug string, params BranchListParams) (*BranchPage, error) {

	q := url.Values{}
	q.Add("base", params.Base)
	q.Add("details", fmt.Sprintf("%v", params.Details))
	q.Add("filterText", params.FilterText)
	q.Add("orderBy", params.OrderBy)
	q.Add("boostMatches", fmt.Sprintf("%v", params.BoostMatches))

	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/branches?%s",
		projectKey,
		repositorySlug,
		q.Encode(),
	)
	_, response, err := c.get(urlPath)
	if err != nil {
		return nil, err
	}
	var branchPage BranchPage
	err = json.Unmarshal(response, &branchPage)
	if err != nil {
		return nil, err
	}
	return &branchPage, nil
}
