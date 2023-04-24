package nexus

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/opendevstack/pipeline/pkg/logging"
	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

// Client represents a Nexus client, wrapping github.com/sonatype-nexus-community/gonexus/rm.RM
type Client struct {
	rm           nexusrm.RM
	httpClient   *http.Client
	clientConfig *ClientConfig
	baseURL      *url.URL
}

// ClientConfig configures a Nexus client.
type ClientConfig struct {
	// Nexus username.
	Username string
	// Password of Nexus user.
	Password string
	// URL of Nexus instance.
	BaseURL string
	// Logger is the logger to send logging messages to.
	Logger logging.LeveledLoggerInterface
	// Timeout of HTTP client used to download assets.
	Timeout time.Duration
	// HTTP client used to download assets.
	HTTPClient *http.Client
	// Certificate to use for HTTP communication.
	CertFile string
}

type ClientInterface interface {
	Download(url, outfile string) (int64, error)
	Search(repository, group string) ([]string, error)
	Upload(repository, group, file string) (string, error)
}

// NewClient initializes a Nexus client.
func NewClient(clientConfig *ClientConfig) (*Client, error) {
	baseURL, err := url.Parse(clientConfig.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	rm, err := nexusrm.New(
		baseURL.String(),
		clientConfig.Username,
		clientConfig.Password,
	)
	if clientConfig.CertFile != "" {
		rm.SetCertFile(clientConfig.CertFile)
	}
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
	return &Client{
		rm:           rm,
		clientConfig: clientConfig,
		httpClient:   httpClient,
		baseURL:      baseURL,
	}, nil
}

// URL returns the Nexus instance URL targeted by this client.
func (c *Client) URL() string {
	return c.baseURL.String()
}

// Username returns the username used by this client.
func (c *Client) Username() string {
	return c.clientConfig.Username
}

// ArtifactGroupBase returns the group base in which aritfacts are stored for
// the given ODS pipeline context.
func ArtifactGroupBase(project, repository, gitCommitSHA string) string {
	return fmt.Sprintf("/%s/%s/%s", project, repository, gitCommitSHA)
}

// ArtifactGroup returns the group in which aritfacts are stored for the given
// ODS pipeline context and the subdir.
func ArtifactGroup(project, repository, gitCommitSHA, subdir string) string {
	return ArtifactGroupBase(project, repository, gitCommitSHA) + "/" + subdir
}

func (c *Client) logger() logging.LeveledLoggerInterface {
	return c.clientConfig.Logger
}

func (c *Client) basicAuth() string {
	auth := c.clientConfig.Username + ":" + c.clientConfig.Password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
