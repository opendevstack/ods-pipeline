package sonar

import (
	"encoding/json"
	"fmt"
)

const (
	TaskStatusInProgress = "IN_PROGRESS"
	TaskStatusPending    = "PENDING"
	TaskStatusSuccess    = "SUCCESS"
	TaskStatusFailed     = "FAILED"
)

type ComputeEngineTask struct {
	Organization       string `json:"organization"`
	ID                 string `json:"id"`
	Type               string `json:"type"`
	ComponentID        string `json:"componentId"`
	ComponentKey       string `json:"componentKey"`
	ComponentName      string `json:"componentName"`
	ComponentQualifier string `json:"componentQualifier"`
	AnalysisID         string `json:"analysisId"`
	Status             string `json:"status"`
	SubmittedAt        string `json:"submittedAt"`
	StartedAt          string `json:"startedAt"`
	ExecutedAt         string `json:"executedAt"`
	ExecutionTimeMs    int    `json:"executionTimeMs"`
	ErrorMessage       string `json:"errorMessage"`
	Logs               bool   `json:"logs"`
	HasErrorStacktrace bool   `json:"hasErrorStacktrace"`
	ErrorStacktrace    string `json:"errorStacktrace"`
	ScannerContext     string `json:"scannerContext"`
	HasScannerContext  bool   `json:"hasScannerContext"`
}

type computeEngineTaskResponse struct {
	Task *ComputeEngineTask `json:"task"`
}

type ComputeEngineTaskGetParams struct {
	AdditionalFields string `json:"additionalFields"`
	ID               string `json:"id"`
}

func (c *Client) ComputeEngineTaskGet(p ComputeEngineTaskGetParams) (*ComputeEngineTask, error) {
	urlPath := fmt.Sprintf("/api/ce/task?id=%s", p.ID)
	statusCode, response, err := c.get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("request returned err: %w", err)
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("request returned unexpected response code: %d, body: %s", statusCode, string(response))
	}
	var cetr *computeEngineTaskResponse
	err = json.Unmarshal(response, &cetr)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %w", err)
	}
	return cetr.Task, nil
}
