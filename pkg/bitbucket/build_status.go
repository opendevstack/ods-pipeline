package bitbucket

import (
	"encoding/json"
	"fmt"
)

type BuildStatusPostPayload struct {
	State       string `json:"state"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// BuildStatusPost associates a build status with a commit.
// The state, the key and the url are mandatory. The name and description fields are optional.
// All fields (mandatory or optional) are limited to 255 characters, except for the url, which is limited to 450 characters.
// Supported values for the state are SUCCESSFUL, FAILED and INPROGRESS.
// The authenticated user must have LICENSED permission or higher to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-build-rest.html#idp6
func (c *Client) BuildStatusPost(gitCommit string, payload BuildStatusPostPayload) error {
	urlPath := fmt.Sprintf("/rest/build-status/1.0/commits/%s", gitCommit)
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	statusCode, response, err := c.post(urlPath, b)
	if err != nil {
		return fmt.Errorf("request returned error: %w", err)
	}
	if statusCode != 204 {
		return fmt.Errorf("request returned unexpected response code: %d, body: %s", statusCode, string(response))
	}
	return nil
}
