package bitbucket

import (
	"encoding/json"
	"fmt"
)

type Tag struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	Hash            string `json:"hash"`
}

type TagCreatePayload struct {
	Message    string `json:"message"`
	Name       string `json:"name"`
	Force      bool   `json:"force"`
	StartPoint string `json:"startPoint"`
	Type       string `json:"type"`
}

// TagCreate creates a tag in the specified repository.
// The authenticated user must have an effective REPO_WRITE permission to call this resource.
//
//'LIGHTWEIGHT' and 'ANNOTATED' are the two type of tags that can be created.
// The 'startPoint' can either be a ref or a 'commit'.
//
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-git-rest.html
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
	if statusCode != 201 {
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
