package tekton

import (
	"context"
	"errors"
	"fmt"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestClient returns mocked pipelines.
type TestClient struct {
	// PipelineRuns is the pool of pipeline runs which can be retrieved.
	PipelineRuns []*tekton.PipelineRun
	// FailCreatePipelineRun lets pipeline run creation fail.
	FailCreatePipelineRun bool
	// CreatedPipelineRuns is a slice of created pipeline run names.
	CreatedPipelineRuns []string
	// FailUpdatePipelineRun lets pipeline run update fail.
	FailUpdatePipelineRun bool
	// UpdatedPipelineRuns is a slice of updated pipeline run names.
	UpdatedPipelineRuns []string
	// FailDeletePipelineRun lets pipeline run deletion fail.
	FailDeletePipelineRun bool
	// DeletedPipelineRuns is a slice of deleted pipeline run names.
	DeletedPipelineRuns []string
}

func (c *TestClient) ListPipelineRuns(ctxt context.Context, options metav1.ListOptions) (*tekton.PipelineRunList, error) {
	items := []tekton.PipelineRun{}
	for _, pr := range c.PipelineRuns {
		items = append(items, *pr)
	}
	return &tekton.PipelineRunList{Items: items}, nil
}

func (c *TestClient) GetPipelineRun(ctxt context.Context, name string, options metav1.GetOptions) (*tekton.PipelineRun, error) {
	for _, p := range c.PipelineRuns {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("pipeline run %s not found", name)
}

func (c *TestClient) CreatePipelineRun(ctxt context.Context, pipeline *tekton.PipelineRun, options metav1.CreateOptions) (*tekton.PipelineRun, error) {
	c.CreatedPipelineRuns = append(c.CreatedPipelineRuns, pipeline.Name)
	if c.FailCreatePipelineRun {
		return nil, errors.New("creation error")
	}
	return pipeline, nil
}

func (c *TestClient) UpdatePipelineRun(ctxt context.Context, pipeline *tekton.PipelineRun, options metav1.UpdateOptions) (*tekton.PipelineRun, error) {
	c.UpdatedPipelineRuns = append(c.UpdatedPipelineRuns, pipeline.Name)
	if c.FailUpdatePipelineRun {
		return nil, errors.New("update error")
	}
	return pipeline, nil
}

func (c *TestClient) DeletePipelineRun(ctxt context.Context, name string, options metav1.DeleteOptions) error {
	c.DeletedPipelineRuns = append(c.DeletedPipelineRuns, name)
	if c.FailDeletePipelineRun {
		return errors.New("delete error")
	}
	return nil
}
