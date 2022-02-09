package tekton

import (
	"context"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientPipelineRunInterface interface {
	ListPipelineRuns(ctxt context.Context, options metav1.ListOptions) (*tekton.PipelineRunList, error)
	GetPipelineRun(ctxt context.Context, name string, options metav1.GetOptions) (*tekton.PipelineRun, error)
	CreatePipelineRun(ctxt context.Context, pipeline *tekton.PipelineRun, options metav1.CreateOptions) (*tekton.PipelineRun, error)
	UpdatePipelineRun(ctxt context.Context, pipeline *tekton.PipelineRun, options metav1.UpdateOptions) (*tekton.PipelineRun, error)
	DeletePipelineRun(ctxt context.Context, name string, options metav1.DeleteOptions) error
}

func (c *Client) ListPipelineRuns(ctxt context.Context, options metav1.ListOptions) (*tekton.PipelineRunList, error) {
	c.logger().Debugf("Get pipeline runs")
	return c.pipelineRunsClient().List(ctxt, options)
}

func (c *Client) GetPipelineRun(ctxt context.Context, name string, options metav1.GetOptions) (*tekton.PipelineRun, error) {
	c.logger().Debugf("Get pipeline run %s", name)
	return c.pipelineRunsClient().Get(ctxt, name, options)
}

func (c *Client) CreatePipelineRun(ctxt context.Context, pipeline *tekton.PipelineRun, options metav1.CreateOptions) (*tekton.PipelineRun, error) {
	c.logger().Debugf("Create pipeline run %s", pipeline.Name)
	return c.pipelineRunsClient().Create(ctxt, pipeline, options)
}

func (c *Client) UpdatePipelineRun(ctxt context.Context, pipeline *tekton.PipelineRun, options metav1.UpdateOptions) (*tekton.PipelineRun, error) {
	c.logger().Debugf("Update pipeline run %s", pipeline.Name)
	return c.pipelineRunsClient().Update(ctxt, pipeline, options)
}

func (c *Client) DeletePipelineRun(ctxt context.Context, name string, options metav1.DeleteOptions) error {
	c.logger().Debugf("Delete pipeline run %s", name)
	return c.pipelineRunsClient().Delete(ctxt, name, options)
}
