package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Tag struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	Hash            string `json:"hash"`
}

type TagPage struct {
	Size       int   `json:"size"`
	Limit      int   `json:"limit"`
	IsLastPage bool  `json:"isLastPage"`
	Values     []Tag `json:"values"`
	Start      int   `json:"start"`
}

type TagCreatePayload struct {
	Message    string `json:"message"`
	Name       string `json:"name"`
	Force      bool   `json:"force"`
	StartPoint string `json:"startPoint"`
	Type       string `json:"type"`
}

type TagListParams struct {
	// FilterText is the the text to match on. The match seems to be a prefix match.
	FilterText string `json:"filterText"`
	// OrderBy determines ordering of refs.
	// Either ALPHABETICAL (by name) or MODIFICATION (last updated).
	OrderBy string `json:"orderBy"`
}

// TagList retrieves the tags matching the supplied filterText param.
// The authenticated user must have REPO_READ permission for the context repository to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp396
func (c *Client) TagList(projectKey string, repositorySlug string, params TagListParams) (*TagPage, error) {

	q := url.Values{}
	q.Add("filterText", params.FilterText)
	q.Add("orderBy", params.OrderBy)

	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/tags?%s",
		projectKey,
		repositorySlug,
		q.Encode(),
	)
	_, response, err := c.get(urlPath)
	if err != nil {
		return nil, err
	}
	var tagPage TagPage
	err = json.Unmarshal(response, &tagPage)
	if err != nil {
		return nil, err
	}
	return &tagPage, nil
}

// TagCreate creates a tag in the specified repository.
// The authenticated user must have an effective REPO_WRITE permission to call this resource.
//
//'LIGHTWEIGHT' and 'ANNOTATED' are the two type of tags that can be created.
// The 'startPoint' can either be a ref or a 'commit'.
//
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp395
func (c *Client) TagCreate(projectKey string, repositorySlug string, payload TagCreatePayload) (*Tag, error) {
	urlPath := fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/tags", projectKey, repositorySlug)
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	statusCode, response, err := c.post(urlPath, b)
	if err != nil {
		return nil, fmt.Errorf("request returned error: %w", err)
	}
	// This endpoint returns 200 based on the documentation and testing. This is
	// contrary to other endpoints which return 201. Therefore we allow both,
	// just in case Atlassian changes their mind in the future to use a proper
	// status code also for this endpoint.
	if statusCode != 201 && statusCode != 200 {
		return nil, fmt.Errorf("request returned unexpected response code: %d, body: %s", statusCode, string(response))
	}
	var tag Tag
	err = json.Unmarshal(response, &tag)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal response: %w. status code: %d, body: %s", err, statusCode, string(response),
		)
	}
	return &tag, nil
}
