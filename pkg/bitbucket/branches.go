package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Branch struct {
	ID              string `json:"id"`
	DisplayId       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	IsDefault       bool   `json:"isDefault"`
}

type BranchPage struct {
	Size       int  `json:"size"`
	Limit      int  `json:"limit"`
	IsLastPage bool `json:"isLastPage"`
	Values     []Branch
	Start      int `json:"start"`
}

type BranchListParams struct {
	Base         string `json:"base"`
	Details      bool   `json:"details"`
	FilterText   string `json:"filterText"`
	OrderBy      string `json:"orderBy"`
	BoostMatches bool   `json:"boostMatches"`
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