package bitbucket

import (
	"encoding/json"
	"fmt"
)

type Repo struct {
	Slug          string  `json:"slug"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Hierarchyid   string  `json:"hierarchyId"`
	SCMID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       Project `json:"project"`
	Public        bool    `json:"public"`
	Links         struct {
		Clone []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

type RepoCreatePayload struct {
	SCMID         string `json:"scmId"`
	Name          string `json:"name"`
	Forkable      bool   `json:"forkable"`
	DefaultBranch string `json:"defaultBranch"`
}

// RepoCreate creates a new repository. Requires an existing project in which this repository will be created. The only parameters which will be used are name and scmId.
// The authenticated user must have PROJECT_ADMIN permission for the context project to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp174
func (c *Client) RepoCreate(projectKey string, payload RepoCreatePayload) (*Repo, error) {
	urlPath := fmt.Sprintf("/rest/api/1.0/projects/%s/repos", projectKey)
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	statusCode, response, err := c.post(urlPath, b)
	if err != nil {
		return nil, fmt.Errorf("request returned error: %w", err)
	}
	var repo Repo
	err = json.Unmarshal(response, &repo)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal response: %w. status code: %d, body: %s", err, statusCode, string(response),
		)
	}
	return &repo, nil
}
