package sonar

import (
	"fmt"
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
	Timeout       time.Duration
	APIToken      string
	HTTPClient    *http.Client
	MaxRetries    int
	BaseURL       string
	ServerEdition string
	Debug         bool
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
		clientConfig.Logger = &logging.LeveledLogger{Level: logging.LevelError}
	}
	if len(clientConfig.ServerEdition) == 0 {
		clientConfig.ServerEdition = "community"
	}
	return &Client{
		httpClient:   httpClient,
		clientConfig: clientConfig,
	}
}

func (c *Client) get(urlPath string) (int, []byte, error) {
	u := c.clientConfig.BaseURL + urlPath
	c.clientConfig.Logger.Debugf("GET %s", u)
	req, err := http.NewRequest("GET", u, nil)
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
