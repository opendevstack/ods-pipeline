package sonar

import (
	"encoding/json"
	"fmt"
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

func (c *Client) QualityGateGet(p QualityGateGetParams) (QualityGate, error) {
	urlPath := fmt.Sprintf("/api/qualitygates/project_status?projectKey=%s", p.Project)
	_, response, err := c.get(urlPath)
	if err != nil {
		return QualityGate{ProjectStatus: QualityGateProjectStatus{Status: "UNKNOWN"}}, nil
	}
	var qg QualityGate
	err = json.Unmarshal(response, &qg)
	if err != nil {
		return qg, err
	}
	return qg, nil
}
