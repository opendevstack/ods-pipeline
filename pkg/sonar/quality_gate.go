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
	Status            string `json:"status"`
	IgnoredConditions bool   `json:"ignoredConditions"`
}

type QualityGateGetParams struct {
	Project string `json:"project"`
}

func (c *Client) QualityGateGet(p QualityGateGetParams) (*QualityGate, error) {
	urlPath := fmt.Sprintf("/api/qualitygates/project_status?projectKey=%s", p.Project)
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
