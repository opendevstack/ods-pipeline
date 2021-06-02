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
