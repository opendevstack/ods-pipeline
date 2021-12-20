package tekton

import (
	"context"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientPipelineInterface interface {
	GetPipeline(ctxt context.Context, name string, options metav1.GetOptions) (*tekton.Pipeline, error)
	CreatePipeline(ctxt context.Context, pipeline *tekton.Pipeline, options metav1.CreateOptions) (*tekton.Pipeline, error)
	UpdatePipeline(ctxt context.Context, pipeline *tekton.Pipeline, options metav1.UpdateOptions) (*tekton.Pipeline, error)
}

func (c *Client) GetPipeline(ctxt context.Context, name string, options metav1.GetOptions) (*tekton.Pipeline, error) {
	c.logger().Debugf("Get pipeline %s", name)
	return c.pipelinesClient().Get(ctxt, name, options)
}

func (c *Client) CreatePipeline(ctxt context.Context, pipeline *tekton.Pipeline, options metav1.CreateOptions) (*tekton.Pipeline, error) {
	c.logger().Debugf("Create pipeline %s", pipeline.Name)
	return c.pipelinesClient().Create(ctxt, pipeline, options)
}

func (c *Client) UpdatePipeline(ctxt context.Context, pipeline *tekton.Pipeline, options metav1.UpdateOptions) (*tekton.Pipeline, error) {
	c.logger().Debugf("Update pipeline %s", pipeline.Name)
	return c.pipelinesClient().Update(ctxt, pipeline, options)
}
