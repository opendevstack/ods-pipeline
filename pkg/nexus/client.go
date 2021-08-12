package nexus

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/opendevstack/pipeline/pkg/logging"
	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

// Client represents a Nexus client, wrapping github.com/sonatype-nexus-community/gonexus/rm.RM
type Client struct {
	rm           nexusrm.RM
	httpClient   *http.Client
	clientConfig *ClientConfig
}

// ClientConfig configures a Nexus client.
type ClientConfig struct {
	// Nexus username.
	Username string
	// Password of Nexus user.
	Password string
	// URL of Nexus instance.
	BaseURL string
	// Nexus repository name.
	Repository string
	// Logger is the logger to send logging messages to.
	Logger logging.LeveledLoggerInterface
	// Timeout of HTTP client used to download assets.
	Timeout time.Duration
	// HTTP client used to download assets.
	HTTPClient *http.Client
}

// NewClient initializes a Nexus client.
func NewClient(clientConfig *ClientConfig) (*Client, error) {
	rm, err := nexusrm.New(
		clientConfig.BaseURL,
		clientConfig.Username,
		clientConfig.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create nexus client: %w", err)
	}
	httpClient := clientConfig.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if clientConfig.Timeout > 0 {
		httpClient.Timeout = clientConfig.Timeout
	} else {
		httpClient.Timeout = 20 * time.Second
	}
	// Be careful not to pass a variable of type *logging.LeveledLogger
	// holding a nil value. If you pass nil for Logger, make sure it is of
	// logging.LeveledLoggerInterface type.
	if clientConfig.Logger == nil {
		clientConfig.Logger = &logging.LeveledLogger{Level: logging.LevelError}
	}

	return &Client{rm: rm, clientConfig: clientConfig, httpClient: httpClient}, nil
}

// URL returns the Nexus instance URL targeted by this client.
func (c *Client) URL() string {
	return c.clientConfig.BaseURL
}

// Repository returns the Nexus repository targeted by this client.
func (c *Client) Repository() string {
	return c.clientConfig.Repository
}

// Username returns the username used by this client.
func (c *Client) Username() string {
	return c.clientConfig.Username
}

func (c *Client) logger() logging.LeveledLoggerInterface {
	return c.clientConfig.Logger
}

func (c *Client) basicAuth() string {
	auth := c.clientConfig.Username + ":" + c.clientConfig.Password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
