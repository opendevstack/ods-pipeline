package bitbucket

import (
	"encoding/json"
	"fmt"
)

type InsightReport struct {
	Data []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
		Type  string `json:"type"`
	} `json:"data"`
	CreatedDate int    `json:"createdDate"`
	Details     string `json:"details"`
	Key         string `json:"key"`
	Link        string `json:"link"`
	LogoURL     string `json:"logoUrl"`
	Result      string `json:"result"`
	Title       string `json:"title"`
	Reporter    string `json:"reporter"`
}

type InsightReportCreatePayload struct {
	Data []struct {
		Title string `json:"title"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"data"`
	Details     string `json:"details"`
	Title       string `json:"title"`
	Reporter    string `json:"reporter"`
	CreatedDate int64  `json:"createdDate"`
	Link        string `json:"link"`
	LogoURL     string `json:"logoUrl"`
	Result      string `json:"result"`
}

// InsightReportCreate creates a new insight report, or replace the existing one if a report already exists for the given repository, commit, and report key. A request to replace an existing report will be rejected if the authenticated user was not the creator of the specified report.
// The report key should be a unique string chosen by the reporter and should be unique enough not to potentially clash with report keys from other reporters. We recommend using reverse DNS namespacing or a similar standard to ensure that collision is avoided.
//
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-code-insights-rest.html#idp9
// TODO: Must be PUT
func (c *Client) InsightReportCreate(projectKey, repositorySlug, commitID, key string, payload InsightReportCreatePayload) (*InsightReport, error) {
	urlPath := fmt.Sprintf(
		"/rest/insights/1.0/projects/%s/repos/%s/commits/%s/reports/%s",
		projectKey,
		repositorySlug,
		commitID,
		key,
	)
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	statusCode, response, err := c.post(urlPath, b)
	if err != nil {
		return nil, fmt.Errorf("request returned error: %w", err)
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("request returned unexpected response code: %d, body: %s", statusCode, string(response))
	}
	var insightReport InsightReport
	err = json.Unmarshal(response, &insightReport)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal response: %w. status code: %d, body: %s", err, statusCode, string(response),
		)
	}
	return &insightReport, nil
}
