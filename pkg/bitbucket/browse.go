package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"
)

type BrowseUpdateParams struct {
	Branch         string
	Message        string
	SourceCommitId string
	Content        io.Reader
}

// BrowseUpdate update the content of path, on the given repository and branch.
// This resource accepts PUT multipart form data, containing the file in a form-field named content.
// An example curl request to update 'README.md' would be:
//  curl -X PUT -u username:password -F content=@README.md  -F 'message=Updated using file-edit REST API'
//  -F branch=master -F  sourceCommitId=5636641a50b
//   http://example.com/rest/api/latest/projects/PROJECT_1/repos/repo_1/browse/README.md
//
// branch: the branch on which the path should be modified or created
// content: the full content of the file at path
// message: the message associated with this change, to be used as the commit message. Or null if the default message should be used.
// sourceCommitId: the commit ID of the file before it was edited, used to identify if content has changed. Or null if this is a new file
//
// The file can be updated or created on a new branch. In this case, the sourceBranch parameter should be provided to identify the starting point for the new branch and the branch parameter identifies the branch to create the new commit on.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp217
func (c *Client) BrowseUpdate(project, repository, path string, params BrowseUpdateParams) (*Commit, error) {
	urlPath := fmt.Sprintf(
		"/rest/api/1.0/projects/%s/repos/%s/browse/%s",
		project,
		repository,
		path,
	)
	paramsMap := map[string]string{
		"branch":         params.Branch,
		"message":        params.Message,
		"sourceCommitId": params.SourceCommitId,
	}
	statusCode, response, err := c.upload(urlPath, paramsMap, path, params.Content)
	if err != nil {
		return nil, fmt.Errorf("could not upload file: %w", err)
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("request returned unexpected status code: %d, body: %s", statusCode, string(response))
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
