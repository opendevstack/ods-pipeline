package bitbucket

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
}

func NewClient(clientConfig *ClientConfig) *Client {
	httpClient := clientConfig.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Client{
		httpClient:   httpClient,
		clientConfig: clientConfig,
	}
}

func (c *Client) get(urlPath string) (int, []byte, error) {
	req, err := http.NewRequest("GET", c.clientConfig.BaseURL+urlPath, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("could not create request: %s", err)
	}

	res, err := c.do(req)
	if err != nil {
		return 500, nil, fmt.Errorf("got error %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	return res.StatusCode, body, err
}

func (c *Client) post(urlPath string, payload []byte) (int, []byte, error) {
	req, err := http.NewRequest("POST", c.clientConfig.BaseURL+urlPath, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, fmt.Errorf("could not create request: %s", err)
	}

	res, err := c.do(req)
	if err != nil {
		return 500, nil, fmt.Errorf("got error %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	return res.StatusCode, body, err
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.clientConfig.APIToken)
	return c.httpClient.Do(req)
}
