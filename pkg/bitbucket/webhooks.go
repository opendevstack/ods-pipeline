package bitbucket

import (
	"encoding/json"
	"fmt"
)

type Webhook struct {
	ID            int                  `json:"id"`
	Name          string               `json:"name"`
	CreatedDate   int64                `json:"createdDate"`
	UpdatedDate   int64                `json:"updatedDate"`
	Events        []string             `json:"events"`
	Configuration WebhookConfiguration `json:"configuration"`
	URL           string               `json:"url"`
	Active        bool                 `json:"active"`
}

type WebhookConfiguration struct {
	Secret string `json:"secret"`
}

type WebhookCreatePayload struct {
	Name          string               `json:"name"`
	Events        []string             `json:"events"`
	Configuration WebhookConfiguration `json:"configuration"`
	URL           string               `json:"url"`
	Active        bool                 `json:"active"`
}

func (c *Client) WebhookCreate(projectKey, repositorySlug string, payload WebhookCreatePayload) (*Webhook, error) {
	urlPath := fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/webhooks", projectKey, repositorySlug)
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	statusCode, response, err := c.post(urlPath, b)
	if err != nil {
		return nil, fmt.Errorf("request returned error: %w", err)
	}
	if statusCode != 201 {
		return nil, fmt.Errorf("request returned unexpected response code: %d, body: %s", statusCode, string(response))
	}
	var webhook Webhook
	err = json.Unmarshal(response, &webhook)
	if err != nil {
		return nil, wrapUnmarshalError(err, statusCode, response)
	}
	return &webhook, nil
}
