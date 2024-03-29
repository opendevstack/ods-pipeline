package bitbucket

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/opendevstack/ods-pipeline/pkg/logging"
)

// Loosely based on https://github.com/brandur/wanikaniapi.
type Client struct {
	httpClient   *http.Client
	clientConfig *ClientConfig
	baseURL      *url.URL
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

func NewClient(clientConfig *ClientConfig) (*Client, error) {
	httpClient := clientConfig.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	// Never follow redirects. Some endpoints (e.g. the one accessed by RawGet)
	// redirect to the login page when authentication is not successful. This
	// behaviour would lead to misleading errors, see e.g.
	// https://github.com/opendevstack/ods-pipeline/issues/254.
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	if clientConfig.Timeout > 0 {
		httpClient.Timeout = clientConfig.Timeout
	} else {
		httpClient.Timeout = 20 * time.Second
	}
	if clientConfig.Logger == nil {
		clientConfig.Logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}
	baseURL, err := url.Parse(clientConfig.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	return &Client{
		httpClient:   httpClient,
		clientConfig: clientConfig,
		baseURL:      baseURL,
	}, nil
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
	u, err := c.baseURL.Parse(urlPath)
	if err != nil {
		return 0, nil, fmt.Errorf("parse URL path: %w", err)
	}
	c.logger().Debugf("%s %s", method, u)
	var requestBody io.Reader
	if payload != nil {
		requestBody = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(method, u.String(), requestBody)
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

	responseBody, err := io.ReadAll(res.Body)
	return res.StatusCode, responseBody, err
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.clientConfig.APIToken)
	return c.httpClient.Do(req)
}

// wrapUnmarshalError wraps err and includes statusCode/response.
func wrapUnmarshalError(err error, statusCode int, response []byte) error {
	return fmt.Errorf("unmarshal: %w. status code: %d, body: %s", err, statusCode, string(response))
}

// fmtStatusCodeError returns an error containing statusCode/response.
func fmtStatusCodeError(statusCode int, response []byte) error {
	return fmt.Errorf("status code: %d, body: %s", statusCode, string(response))
}
