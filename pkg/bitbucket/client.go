package bitbucket

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/opendevstack/pipeline/pkg/logging"
)

// Loosely based on https://github.com/brandur/wanikaniapi.
type Client struct {
	httpClient   *http.Client
	clientConfig *ClientConfig
}

type ClientConfig struct {
	Timeout    time.Duration
	APIToken   string
	HTTPClient *http.Client
	MaxRetries int
	BaseURL    string
	// Logger is the logger to send logging messages to.
	Logger logging.LeveledLoggerInterface
}

func NewClient(clientConfig *ClientConfig) *Client {
	httpClient := clientConfig.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if clientConfig.Timeout > 0 {
		httpClient.Timeout = clientConfig.Timeout
	} else {
		httpClient.Timeout = 20 * time.Second
	}
	if clientConfig.Logger == nil {
		clientConfig.Logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}
	return &Client{
		httpClient:   httpClient,
		clientConfig: clientConfig,
	}
}

func (c *Client) get(urlPath string) (int, []byte, error) {
	return c.createRequest("GET", urlPath, nil)
}

func (c *Client) post(urlPath string, payload []byte) (int, []byte, error) {
	return c.createRequest("POST", urlPath, payload)
}

func (c *Client) put(urlPath string, payload []byte) (int, []byte, error) {
	return c.createRequest("PUT", urlPath, payload)
}

func (c *Client) createRequest(method, urlPath string, payload []byte) (int, []byte, error) {
	u := c.clientConfig.BaseURL + urlPath
	c.logger().Debugf("%s %s", method, u)
	var requestBody io.Reader
	if payload != nil {
		requestBody = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(method, u, requestBody)
	if err != nil {
		return 0, nil, fmt.Errorf("could not create request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req)
}

func (c *Client) logger() logging.LeveledLoggerInterface {
	return c.clientConfig.Logger
}

func (c *Client) doRequest(req *http.Request) (int, []byte, error) {
	res, err := c.do(req)
	if err != nil {
		return 500, nil, fmt.Errorf("got error %s", err)
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	return res.StatusCode, responseBody, err
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.clientConfig.APIToken)
	return c.httpClient.Do(req)
}
