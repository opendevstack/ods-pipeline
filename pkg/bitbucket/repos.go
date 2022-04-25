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

type RepoPage struct {
	Size       int    `json:"size"`
	Limit      int    `json:"limit"`
	IsLastPage bool   `json:"isLastPage"`
	Values     []Repo `json:"values"`
	Start      int    `json:"start"`
}

type RepoCreatePayload struct {
	SCMID         string `json:"scmId"`
	Name          string `json:"name"`
	Forkable      bool   `json:"forkable"`
	DefaultBranch string `json:"defaultBranch"`
}

type RepoClientInterface interface {
	RepoList(projectKey string) (*RepoPage, error)
	RepoCreate(projectKey string, payload RepoCreatePayload) (*Repo, error)
}

// RepoList retrieves repositories from the project corresponding to the supplied projectKey.
// The authenticated user must have REPO_READ permission for the context repository to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp175
func (c *Client) RepoList(projectKey string) (*RepoPage, error) {

	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos",
		projectKey,
	)
	statusCode, response, err := c.get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("retrieve %s: %w", urlPath, err)
	}
	if statusCode != 200 {
		return nil, fmtStatusCodeError(statusCode, response)
	}
	var repoPage RepoPage
	err = json.Unmarshal(response, &repoPage)
	if err != nil {
		return nil, wrapUnmarshalError(err, statusCode, response)
	}
	return &repoPage, nil
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
		return nil, fmt.Errorf("create %s: %w", projectKey, err)
	}
	if statusCode != 201 {
		return nil, fmtStatusCodeError(statusCode, response)
	}
	var repo Repo
	err = json.Unmarshal(response, &repo)
	if err != nil {
		return nil, wrapUnmarshalError(err, statusCode, response)
	}
	return &repo, nil
}
