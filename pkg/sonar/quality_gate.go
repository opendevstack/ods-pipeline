package sonar

import (
	"encoding/json"
	"fmt"
)

const (
	QualityGateStatusOk    = "OK"
	QualityGateStatusWarn  = "WARN"
	QualityGateStatusError = "ERROR"
	QualityGateStatusNone  = "NONE"
)

type QualityGate struct {
	ProjectStatus QualityGateProjectStatus `json:"projectStatus"`
}

type QualityGateProjectStatus struct {
	Status            string                 `json:"status"`
	IgnoredConditions bool                   `json:"ignoredConditions"`
	Conditions        []QualityGateCondition `json:"conditions"`
	Periods           []QualityGatePeriod    `json:"periods"`
}

type QualityGateCondition struct {
	Status         string `json:"status"`
	MetricKey      string `json:"metricKey"`
	Comparator     string `json:"comparator"`
	PeriodIndex    int    `json:"periodIndex"`
	ErrorThreshold string `json:"errorThreshold,omitempty"`
	ActualValue    string `json:"actualValue"`
}

type QualityGatePeriod struct {
	Index     int    `json:"index"`
	Mode      string `json:"mode"`
	Date      string `json:"date"`
	Parameter string `json:"parameter"`
}

type QualityGateGetParams struct {
	ProjectKey  string
	Branch      string
	PullRequest string
}

func (c *Client) QualityGateGet(p QualityGateGetParams) (*QualityGate, error) {
	urlPath := "/api/qualitygates/project_status?projectKey=" + p.ProjectKey
	if p.PullRequest != "" && p.PullRequest != "0" {
		urlPath = urlPath + "&pullRequest=" + p.PullRequest
	} else if p.Branch != "" {
		urlPath = urlPath + "&branch=" + p.Branch
	}
	statusCode, response, err := c.get(urlPath)
	if err != nil {
		return &QualityGate{ProjectStatus: QualityGateProjectStatus{Status: QualityGateStatusNone}}, nil
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("request returned unexpected response code: %d, body: %s", statusCode, string(response))
	}
	var qg *QualityGate
	err = json.Unmarshal(response, &qg)
	if err != nil {
		return qg, err
	}
	return qg, nil
}
