package bitbucket

import (
	"encoding/json"
	"fmt"
)

type Project struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

type ProjectCreatePayload struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

// ProjectCreate creates a new project.
// To include a custom avatar for the project, the project definition should contain an additional attribute with the key avatar and the value a data URI containing Base64-encoded image data. The URI should be in the following format:
//      data:(content type, e.g. image/png);base64,(data)
// If the data is not Base64-encoded, or if a character set is defined in the URI, or the URI is otherwise invalid, project creation will fail.
// The authenticated user must have PROJECT_CREATE permission to call this resource.
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-rest.html#idp148
func (c *Client) ProjectCreate(payload ProjectCreatePayload) (*Project, error) {
	urlPath := "/rest/api/1.0/projects"
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	statusCode, response, err := c.post(urlPath, b)
	if err != nil {
		return nil, fmt.Errorf("request returned error: %w", err)
	}
	if statusCode != 201 {
		return nil, fmt.Errorf("request returned unexpected status code: %d, body: %s", statusCode, string(response))
	}
	var project Project
	err = json.Unmarshal(response, &project)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal response: %w. status code: %d, body: %s", err, statusCode, string(response),
		)
	}
	return &project, nil
}
