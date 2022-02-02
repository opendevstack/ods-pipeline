package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NotificationConfigMap   = "ods-notification"
	UrlProperty             = "url"
	MethodProperty          = "method"
	ContentTypeProperty     = "contentType"
	RequestTemplateProperty = "requestTemplate"
	NotifyOnStatusProperty  = "notifyOnStatus"
	EnabledProperty         = "enabled"
)

type Client struct {
	clientConfig     ClientConfig
	httpClient       *http.Client
	kubernetesClient kubernetes.ClientInterface
}

type ClientConfig struct {
	Namespace  string
	Logger     logging.LeveledLoggerInterface
	HTTPClient *http.Client
}

type PipelineRunResult struct {
	PipelineRunName string
	PipelineRunURL  string
	OverallStatus   string
	ODSContext      *pipelinectxt.ODSContext
}

type NotificationConfig struct {
	enabled        bool
	url            string
	method         string
	contentType    string
	notifyOnStatus []string
	template       *template.Template
}

func NewClient(config ClientConfig, kubernetesClient kubernetes.ClientInterface) (*Client, error) {
	if config.Logger == nil {
		config.Logger = &logging.LeveledLogger{Level: logging.LevelError}
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if httpClient.Timeout == 0 {
		httpClient.Timeout = 20 * time.Second
	}

	return &Client{
		clientConfig:     config,
		kubernetesClient: kubernetesClient,
		httpClient:       httpClient,
	}, nil
}

func (c Client) ReadNotificationConfig(ctxt context.Context) (*NotificationConfig, error) {
	cm, err := c.kubernetesClient.GetConfigMap(ctxt, NotificationConfigMap, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load %s ConfigMap: %v", NotificationConfigMap, err)
	}

	enabledPropValue, ok := cm.Data[EnabledProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", NotificationConfigMap, EnabledProperty)
	}

	enabled, err := strconv.ParseBool(enabledPropValue)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %s to bool", enabledPropValue)
	}

	url, ok := cm.Data[UrlProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", NotificationConfigMap, UrlProperty)
	}

	method, ok := cm.Data[MethodProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", NotificationConfigMap, MethodProperty)
	}

	contentType, ok := cm.Data[ContentTypeProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", NotificationConfigMap, ContentTypeProperty)
	}

	notifyOnStatus, ok := cm.Data[NotifyOnStatusProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specifiy '%s' property", NotificationConfigMap, NotifyOnStatusProperty)
	}

	decoder := json.NewDecoder(strings.NewReader(notifyOnStatus))
	var notificationStatusValues []string
	err = decoder.Decode(&notificationStatusValues)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("decoding notification status properties failed: %w", err)
	}

	text, ok := cm.Data[RequestTemplateProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", NotificationConfigMap, RequestTemplateProperty)
	}

	requestTemplate, err := template.New("requestTemplate").Parse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse requestTemplate template")
	}

	return &NotificationConfig{
		enabled,
		url,
		method,
		contentType,
		notificationStatusValues,
		requestTemplate,
	}, nil
}

func (c *Client) ShouldNotify(notificationConfig *NotificationConfig, status string) bool {
	if notificationConfig.enabled {
		for _, allowedStatus := range notificationConfig.notifyOnStatus {
			if allowedStatus == status {
				return true
			}
		}
	}
	return false
}

func (c Client) CallWebhook(ctxt context.Context, notificationConfig *NotificationConfig, summary PipelineRunResult) error {
	requestBody := bytes.NewBuffer([]byte{})
	err := notificationConfig.template.Execute(requestBody, summary)
	if err != nil {
		return fmt.Errorf("rendering notification webhook template failed: %v", err)
	}

	req, err := http.NewRequest(notificationConfig.method, notificationConfig.url, requestBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", notificationConfig.contentType)
	req.WithContext(ctxt)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("performing notification webhook request failed: %v", err)
	}
	c.logger().Infof("notification webhook response was: %s", response.Status)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			body = []byte("<could not read body>")
		}
		return fmt.Errorf("notification webhook returned non-2xx status (status: %s, body: %s)",
			response.Status, string(body))
	}
	// we do not fail
	return nil
}

func (c Client) logger() logging.LeveledLoggerInterface {
	return c.clientConfig.Logger
}
