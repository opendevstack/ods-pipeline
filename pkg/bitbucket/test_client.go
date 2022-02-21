package bitbucket

import (
	"errors"
	"fmt"
)

// TestClient returns mocked branches and tags.
type TestClient struct {
	Branches     []Branch
	Tags         []Tag
	Repos        []Repo
	Commits      []Commit
	PullRequests []PullRequest
	// Files contains byte slices for filenames
	Files map[string][]byte
}

func (c *TestClient) BranchList(projectKey string, repositorySlug string, params BranchListParams) (*BranchPage, error) {
	return &BranchPage{
		Values: c.Branches,
	}, nil
}

func (c *TestClient) TagList(projectKey string, repositorySlug string, params TagListParams) (*TagPage, error) {
	return &TagPage{
		Values: c.Tags,
	}, nil
}

func (c *TestClient) TagGet(projectKey string, repositorySlug string, name string) (*Tag, error) {
	for _, t := range c.Tags {
		if t.DisplayID == name {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("no tag %s", name)
}

func (c *TestClient) TagCreate(projectKey string, repositorySlug string, payload TagCreatePayload) (*Tag, error) {
	return nil, errors.New("not implemented")
}

func (c *TestClient) RepoList(projectKey string) (*RepoPage, error) {
	return &RepoPage{
		Values: c.Repos,
	}, nil
}

func (c *TestClient) RepoCreate(projectKey string, payload RepoCreatePayload) (*Repo, error) {
	return nil, errors.New("not implemented")
}

func (c *TestClient) RawGet(project, repository, filename, gitFullRef string) ([]byte, error) {
	if f, ok := c.Files[filename]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("%s not found", filename)
}

func (c *TestClient) CommitList(projectKey string, repositorySlug string, params CommitListParams) (*CommitPage, error) {
	return &CommitPage{Values: c.Commits}, nil
}

func (c *TestClient) CommitGet(projectKey, repositorySlug, commitID string) (*Commit, error) {
	for _, co := range c.Commits {
		if co.ID == commitID {
			return &co, nil
		}
	}
	return nil, fmt.Errorf("no commit %s", commitID)
}

func (c *TestClient) CommitPullRequestList(projectKey, repositorySlug, commitID string) (*PullRequestPage, error) {
	return &PullRequestPage{Values: c.PullRequests}, nil
}
