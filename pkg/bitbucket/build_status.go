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

// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-build-rest.html#idp6
func (c *Client) BuildStatusPost(gitCommit string, payload BuildStatusPostPayload) (int, []byte, error) {
	urlPath := fmt.Sprintf("/rest/build-status/1.0/commits/%s", gitCommit)
	b, err := json.Marshal(payload)
	if err != nil {
		return 0, []byte{}, err
	}
	return c.post(urlPath, b)
}
