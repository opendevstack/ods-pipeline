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
		return 500, nil, fmt.Errorf("%s %s: %w", req.Method, req.URL, err)
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("read %s: %w", req.URL, err)
	}
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

// wrapRequestError wraps an error for a request.
func wrapRequestError(err error) error {
	return fmt.Errorf("request: %w", err)
}
