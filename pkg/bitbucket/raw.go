package bitbucket

import (
	"fmt"
	"net/url"
)

// RawGet retrieves the raw content for a file path at a specified revision.
// The authenticated user must have REPO_READ permission for the specified repository to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp359
func (c *Client) RawGet(project, repository, filename, gitFullRef string) ([]byte, error) {
	urlPath := fmt.Sprintf(
		"/projects/%s/repos/%s/raw/%s?at=%s",
		project,
		repository,
		filename,
		url.QueryEscape(gitFullRef),
	)
	fmt.Println(urlPath)
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
