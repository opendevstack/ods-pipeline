package notification

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

type Client struct {
	clientConfig ClientConfig
	httpClient   *http.Client
}

type ClientConfig struct {
	Namespace          string
	Logger             logging.LeveledLoggerInterface
	NotificationConfig *Config
	HTTPClient         *http.Client
}

type PipelineRunResult struct {
	PipelineRunName string
	PipelineRunURL  string
	OverallStatus   string
	ODSContext      *pipelinectxt.ODSContext
}

func NewClient(config ClientConfig) (*Client, error) {
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
		clientConfig: config,
		httpClient:   httpClient,
	}, nil
}

func (c *Client) ShouldNotify(status string) bool {
	notificationConfig := c.clientConfig.NotificationConfig
	if notificationConfig.Enabled {
		for _, allowedStatus := range notificationConfig.NotifyOnStatus {
			if allowedStatus == status {
				return true
			}
		}
	}
	return false
}

func (c Client) CallWebhook(ctxt context.Context, summary PipelineRunResult) error {
	notificationConfig := c.clientConfig.NotificationConfig
	requestBody := bytes.NewBuffer([]byte{})
	err := notificationConfig.Template.Execute(requestBody, summary)
	if err != nil {
		return fmt.Errorf("rendering notification webhook Template failed: %w", err)
	}

	req, err := http.NewRequest(notificationConfig.Method, notificationConfig.URL, requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", notificationConfig.ContentType)
	req = req.WithContext(ctxt)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("performing notification webhook request failed: %w", err)
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
