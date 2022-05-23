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
	Size        int      `json:"size"`
	Limit       int      `json:"limit"`
	IsLastPage  bool     `json:"isLastPage"`
	Values      []Commit `json:"values"`
	Start       int      `json:"start"`
	AuthorCount int      `json:"authorCount"`
	TotalCount  int      `json:"totalCount"`
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

type CommitClientInterface interface {
	CommitList(projectKey string, repositorySlug string, params CommitListParams) (*CommitPage, error)
	CommitGet(projectKey, repositorySlug, commitID string) (*Commit, error)
	CommitPullRequestList(projectKey, repositorySlug, commitID string) (*PullRequestPage, error)
}

// CommitList retrieves a page of commits from a given starting commit or "between" two commits. If no explicit commit is specified, the tip of the repository's default branch is assumed. commits may be identified by branch or tag name or by ID. A path may be supplied to restrict the returned commits to only those which affect that path.
// The authenticated user must have REPO_READ permission for the specified repository to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp222
func (c *Client) CommitList(projectKey string, repositorySlug string, params CommitListParams) (*CommitPage, error) {

	q := url.Values{}
	q.Add("since", params.Since)
	q.Add("until", params.Until)

	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/commits?%s",
		projectKey,
		repositorySlug,
		q.Encode(),
	)
	statusCode, response, err := c.get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("retrieve %s: %w", urlPath, err)
	}
	if statusCode != 200 {
		return nil, fmtStatusCodeError(statusCode, response)
	}
	var commitPage CommitPage
	err = json.Unmarshal(response, &commitPage)
	if err != nil {
		return nil, wrapUnmarshalError(err, statusCode, response)
	}
	return &commitPage, nil
}

// CommitGet etrieves a single commit identified by its ID. In general, that ID is a SHA1. From 2.11, ref names like "refs/heads/master" are no longer accepted by this resource.
// The authenticated user must have REPO_READ permission for the specified repository to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp224
func (c *Client) CommitGet(projectKey, repositorySlug, commitID string) (*Commit, error) {
	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/commits/%s",
		projectKey,
		repositorySlug,
		commitID,
	)
	statusCode, response, err := c.get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("retrieve %s: %w", urlPath, err)
	}
	if statusCode != 200 {
		return nil, fmtStatusCodeError(statusCode, response)
	}
	var commit Commit
	err = json.Unmarshal(response, &commit)
	if err != nil {
		return nil, wrapUnmarshalError(err, statusCode, response)
	}
	return &commit, nil
}

// CommitPullRequestList retrieves a page of pull requests in the current repository that contain the given commit.
// The user must be authenticated and have access to the specified repository to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp243
func (c *Client) CommitPullRequestList(projectKey, repositorySlug, commitID string) (*PullRequestPage, error) {
	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/commits/%s/pull-requests",
		projectKey,
		repositorySlug,
		commitID,
	)
	statusCode, response, err := c.get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("retrieve %s: %w", urlPath, err)
	}
	if statusCode != 200 {
		return nil, fmtStatusCodeError(statusCode, response)
	}
	var prPage PullRequestPage
	err = json.Unmarshal(response, &prPage)
	if err != nil {
		return nil, wrapUnmarshalError(err, statusCode, response)
	}
	return &prPage, nil
}
