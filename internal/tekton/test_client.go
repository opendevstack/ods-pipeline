package tekton

import (
	"context"
	"errors"
	"fmt"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestClient returns mocked pipelines.
type TestClient struct {
	// Pipelines is the pool of pipelines which can be retrieved.
	Pipelines []*tekton.Pipeline
	// FailCreatePipeline lets pipeline creation fail.
	FailCreatePipeline bool
	// CreatedPipelines is a slice of created pipeline names.
	CreatedPipelines []string
	// FailUpdatePipeline lets pipeline update fail.
	FailUpdatePipeline bool
	// UpdatedPipelines is a slice of updated pipeline names.
	UpdatedPipelines []string
}

func (c *TestClient) GetPipeline(ctxt context.Context, name string, options metav1.GetOptions) (*tekton.Pipeline, error) {
	for _, p := range c.Pipelines {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("pipeline %s not found", name)
}

func (c *TestClient) CreatePipeline(ctxt context.Context, pipeline *tekton.Pipeline, options metav1.CreateOptions) (*tekton.Pipeline, error) {
	c.CreatedPipelines = append(c.CreatedPipelines, pipeline.Name)
	if c.FailCreatePipeline {
		return nil, errors.New("creation error")
	}
	return pipeline, nil
}

func (c *TestClient) UpdatePipeline(ctxt context.Context, pipeline *tekton.Pipeline, options metav1.UpdateOptions) (*tekton.Pipeline, error) {
	c.UpdatedPipelines = append(c.UpdatedPipelines, pipeline.Name)
	if c.FailUpdatePipeline {
		return nil, errors.New("update error")
	}
	return pipeline, nil
}
