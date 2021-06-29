package bitbucket

import (
	"fmt"
	"net/url"
)

// RawGet retrieves the raw content for a file path at a specified revision.
// Update the content of path, on the given repository and branch.
This resource accepts PUT multipart form data, containing the file in a form-field named content.
An example curl request to update 'README.md' would be:
 curl -X PUT -u username:password -F content=@README.md  -F 'message=Updated using file-edit REST API'
 -F branch=master -F  sourceCommitId=5636641a50b
  http://example.com/rest/api/latest/projects/PROJECT_1/repos/repo_1/browse/README.md
 
branch: the branch on which the path should be modified or created
content: the full content of the file at path
message: the message associated with this change, to be used as the commit message. Or null if the default message should be used.
sourceCommitId: the commit ID of the file before it was edited, used to identify if content has changed. Or null if this is a new file
The file can be updated or created on a new branch. In this case, the sourceBranch parameter should be provided to identify the starting point for the new branch and the branch parameter identifies the branch to create the new commit on.
// https://docs.atlassian.com/bitbucket-server/rest/6.4.0/bitbucket-rest.html#idp188
func (c *Client) BrowseGet(project, repository, filename, gitFullRef string) ([]byte, error) {
	urlPath := fmt.Sprintf(
		"/projects/%s/repos/%s/raw/%s?at=%s",
		project,
		repository,
		filename,
		url.QueryEscape(gitFullRef),
	)
	statusCode, body, err := c.get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("could not get file: %w", err)
	}

	switch statusCode {
	case 200:
		return body, nil
	case 404:
		return nil, fmt.Errorf("could not find file '%s' at '%s'", filename, gitFullRef)
	default:
		return nil, fmt.Errorf("unexpected status code %d", statusCode)
	}
}
