package webhook

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	WebhookConfigMap        = "webhook-config"
	UrlProperty             = "url"
	MethodProperty          = "method"
	ContentTypeProperty     = "contentType"
	RequestTemplateProperty = "request"
)

type Client struct {
	clientConfig     ClientConfig
	httpClient       *http.Client
	kubernetesClient kubernetes.ClientInterface
}

type ClientConfig struct {
	Namespace string
	Logger    logging.LeveledLoggerInterface
}

type PipelineRunResult struct {
	PipelineRunURL string
	OverallStatus  string
	ODSContext     *pipelinectxt.ODSContext
}

type webhookConfig struct {
	url         string
	method      string
	contentType string
	template    *template.Template
}

func NewClient(config ClientConfig, kubernetesClient kubernetes.ClientInterface) (*Client, error) {
	if config.Logger == nil {
		config.Logger = &logging.LeveledLogger{Level: logging.LevelError}
	}

	return &Client{
		clientConfig:     config,
		httpClient:       &http.Client{},
		kubernetesClient: kubernetesClient,
	}, nil
}

func (c Client) readWebhookConfig(ctxt context.Context) (*webhookConfig, error) {
	cm, err := c.kubernetesClient.GetConfigMap(ctxt, WebhookConfigMap, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load webhook-config ConfigMap: %v", err)
	}

	url, ok := cm.Data[UrlProperty]
	if !ok {
		return nil, fmt.Errorf("webhook-config doesn't specify '%s' property", UrlProperty)
	}

	method, ok := cm.Data[MethodProperty]
	if !ok {
		return nil, fmt.Errorf("webhook-config doesn't specify '%s' property", MethodProperty)
	}

	contentType, ok := cm.Data[ContentTypeProperty]
	if !ok {
		return nil, fmt.Errorf("webhook-config doesn't specify '%s' property", ContentTypeProperty)
	}

	text, ok := cm.Data[RequestTemplateProperty]
	if !ok {
		return nil, fmt.Errorf("webhook-config doesn't specify '%s' property", RequestTemplateProperty)
	}

	requestTemplate, err := template.New("request-template").Parse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request template")
	}

	return &webhookConfig{
		url,
		method,
		contentType,
		requestTemplate,
	}, nil
}

func (c Client) CallWebhook(ctxt context.Context, summary PipelineRunResult) error {
	config, err := c.readWebhookConfig(ctxt)
	if err != nil {
		return fmt.Errorf("unable to read webhook config: %v", err)
	}

	requestBody := bytes.NewBuffer([]byte{})
	if config.template.Execute(requestBody, summary) != nil {
		return fmt.Errorf("rendering webhook template failed: %v", err)
	}

	req, err := http.NewRequest(config.method, config.url, requestBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", config.contentType)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("performing webhook request failed: %v", err)
	}
	c.logger().Infof("Notification Webhook response was: %w", response.StatusCode)
	// we do not fail
	return nil
}

func (c Client) logger() logging.LeveledLoggerInterface {
	return c.clientConfig.Logger
}
