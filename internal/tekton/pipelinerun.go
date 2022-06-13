package tekton

import (
	"context"
	"fmt"
	"net/url"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientPipelineRunInterface interface {
	ListPipelineRuns(ctxt context.Context, options metav1.ListOptions) (*tekton.PipelineRunList, error)
	GetPipelineRun(ctxt context.Context, name string, options metav1.GetOptions) (*tekton.PipelineRun, error)
	CreatePipelineRun(ctxt context.Context, pipelineRun *tekton.PipelineRun, options metav1.CreateOptions) (*tekton.PipelineRun, error)
	UpdatePipelineRun(ctxt context.Context, pipelineRun *tekton.PipelineRun, options metav1.UpdateOptions) (*tekton.PipelineRun, error)
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

func (c *Client) CreatePipelineRun(ctxt context.Context, pipelineRun *tekton.PipelineRun, options metav1.CreateOptions) (*tekton.PipelineRun, error) {
	c.logger().Debugf("Create pipeline run %s", pipelineRun.Name)
	return c.pipelineRunsClient().Create(ctxt, pipelineRun, options)
}

func (c *Client) UpdatePipelineRun(ctxt context.Context, pipelineRun *tekton.PipelineRun, options metav1.UpdateOptions) (*tekton.PipelineRun, error) {
	c.logger().Debugf("Update pipeline run %s", pipelineRun.Name)
	return c.pipelineRunsClient().Update(ctxt, pipelineRun, options)
}

func (c *Client) DeletePipelineRun(ctxt context.Context, name string, options metav1.DeleteOptions) error {
	c.logger().Debugf("Delete pipeline run %s", name)
	return c.pipelineRunsClient().Delete(ctxt, name, options)
}

// PipelineRunURL returns an URL to the pipeline run given in opts.
func PipelineRunURL(consoleURL, namespace, name string) (string, error) {
	cURL, err := url.Parse(consoleURL)
	if err != nil {
		return "", fmt.Errorf("parse base URL: %w", err)
	}
	cPath := fmt.Sprintf(
		"/k8s/ns/%s/tekton.dev~v1beta1~PipelineRun/%s/",
		namespace,
		name,
	)
	fullURL, err := cURL.Parse(cPath)
	if err != nil {
		return "", fmt.Errorf("parse URL path: %w", err)
	}
	return fullURL.String(), nil
}
