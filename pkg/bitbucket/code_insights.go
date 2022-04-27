package bitbucket

import (
	"encoding/json"
	"fmt"
)

const (
	InsightReportPass = "PASS"
	InsightReportFail = "FAIL"
)

type InsightReport struct {
	Data        []InsightReportData `json:"data"`
	CreatedDate int                 `json:"createdDate"`
	Details     string              `json:"details"`
	Key         string              `json:"key"`
	Link        string              `json:"link"`
	LogoURL     string              `json:"logoUrl"`
	Result      string              `json:"result"`
	Title       string              `json:"title"`
	Reporter    string              `json:"reporter"`
}

type InsightReportCreatePayload struct {
	Data        []InsightReportData `json:"data"`
	Details     string              `json:"details,omitempty"`
	Title       string              `json:"title"`
	Reporter    string              `json:"reporter,omitempty"`
	CreatedDate int64               `json:"createdDate"`
	Link        string              `json:"link,omitempty"`
	LogoURL     string              `json:"logoUrl,omitempty"`
	Result      string              `json:"result,omitempty"`
}

type InsightReportData struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type CodeInsightsClientInterface interface {
	InsightReportCreate(projectKey, repositorySlug, commitID, key string, payload InsightReportCreatePayload) (*InsightReport, error)
}

// InsightReportCreate creates a new insight report, or replace the existing one if a report already exists for the given repository, commit, and report key. A request to replace an existing report will be rejected if the authenticated user was not the creator of the specified report.
// The report key should be a unique string chosen by the reporter and should be unique enough not to potentially clash with report keys from other reporters. We recommend using reverse DNS namespacing or a similar standard to ensure that collision is avoided.
//
// https://docs.atlassian.com/bitbucket-server/rest/7.13.0/bitbucket-code-insights-rest.html#idp9
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
	statusCode, response, err := c.put(urlPath, b)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	if statusCode != 200 {
		c.clientConfig.Logger.Debugf("Request Body:\n%s", string(b))
		return nil, fmtStatusCodeError(statusCode, response)
	}
	var insightReport InsightReport
	err = json.Unmarshal(response, &insightReport)
	if err != nil {
		return nil, wrapUnmarshalError(err, statusCode, response)
	}
	return &insightReport, nil
}
