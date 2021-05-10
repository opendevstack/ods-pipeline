package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Commit struct {
	ID        string `json:"id"`
	DisplayID string `json:"displayId"`
	Author    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
	} `json:"author"`
	AuthorTimestamp int64 `json:"authorTimestamp"`
	Committer       struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
	} `json:"committer"`
	CommitterTimestamp int64  `json:"committerTimestamp"`
	Message            string `json:"message"`
	Parents            []struct {
		ID        string `json:"id"`
		DisplayID string `json:"displayId"`
	} `json:"parents"`
}

type CommitPage struct {
	Size        int  `json:"size"`
	Limit       int  `json:"limit"`
	IsLastPage  bool `json:"isLastPage"`
	Values      []Commit
	Start       int `json:"start"`
	AuthorCount int `json:"authorCount"`
	TotalCount  int `json:"totalCount"`
}

type PullRequestPage struct {
	Size       int  `json:"size"`
	Limit      int  `json:"limit"`
	IsLastPage bool `json:"isLastPage"`
	Values     []PullRequest
	Start      int `json:"start"`
}

type PullRequest struct {
	ID          int    `json:"id"`
	Version     int    `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	Open        bool   `json:"open"`
	Closed      bool   `json:"closed"`
	CreatedDate int    `json:"createdDate"`
	UpdatedDate int    `json:"updatedDate"`
	FromRef     Ref    `json:"fromRef"`
	ToRef       Ref    `json:"toRef"`
	Locked      bool   `json:"locked"`
	Author      struct {
		User     User   `json:"user"`
		Role     string `json:"role"`
		Approved bool   `json:"approved"`
		Status   string `json:"status"`
	} `json:"author"`
	Reviewers []struct {
		User               User   `json:"user"`
		LastReviewedCommit string `json:"lastReviewedCommit"`
		Role               string `json:"role"`
		Approved           bool   `json:"approved"`
		Status             string `json:"status"`
	} `json:"reviewers"`
	Participants []struct {
		User     User   `json:"user"`
		Role     string `json:"role"`
		Approved bool   `json:"approved"`
		Status   string `json:"status"`
	} `json:"participants"`
	Links struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

type Ref struct {
	ID         string `json:"id"`
	Repository struct {
		Slug    string      `json:"slug"`
		Name    interface{} `json:"name"`
		Project struct {
			Key string `json:"key"`
		} `json:"project"`
	} `json:"repository"`
}

type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type CommitListParams struct {
	Since string `json:"since"`
	Until string `json:"until"`
}

func (c *Client) CommitList(projectKey string, repositorySlug string, params CommitListParams) (*CommitPage, error) {

	q := url.Values{}
	q.Add("api_key", "key_from_environment_or_flag")
	q.Add("another_thing", "foo & bar")

	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/commits?%s",
		projectKey,
		repositorySlug,
		q.Encode(),
	)
	_, response, err := c.get(urlPath)
	if err != nil {
		return nil, err
	}
	var commitPage CommitPage
	err = json.Unmarshal(response, &commitPage)
	if err != nil {
		return nil, err
	}
	return &commitPage, nil
}

func (c *Client) CommitGet(projectKey, repositorySlug, commitID string) (*Commit, error) {
	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/commits/%s",
		projectKey,
		repositorySlug,
		commitID,
	)
	statusCode, response, err := c.get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("request returned error: %w", err)
	}
	var commit Commit
	err = json.Unmarshal(response, &commit)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal response: %w. status code: %d, body: %s", err, statusCode, string(response),
		)
	}
	return &commit, nil
}

func (c *Client) CommitPullRequestList(projectKey, repositorySlug, commitID string) (*PullRequestPage, error) {
	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/commits/%s/pull-requests",
		projectKey,
		repositorySlug,
		commitID,
	)
	statusCode, response, err := c.get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("request returned error: %w", err)
	}
	var prPage PullRequestPage
	err = json.Unmarshal(response, &prPage)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal response: %w. status code: %d, body: %s", err, statusCode, string(response),
		)
	}
	return &prPage, nil
}
