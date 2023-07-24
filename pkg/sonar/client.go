package sonar

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/opendevstack/ods-pipeline/pkg/logging"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
)

type ClientInterface interface {
	Scan(sonarProject, branch, commit string, pr *PullRequest, outWriter, errWriter io.Writer) error
	QualityGateGet(p QualityGateGetParams) (*QualityGate, error)
	GenerateReports(sonarProject, author, branch, rootPath, artifactPrefix string) error
	ExtractComputeEngineTaskID(filename string) (string, error)
	ComputeEngineTaskGet(p ComputeEngineTaskGetParams) (*ComputeEngineTask, error)
}

// Loosely based on https://github.com/brandur/wanikaniapi.
type Client struct {
	httpClient   *http.Client
	clientConfig *ClientConfig
	baseURL      *url.URL
}

type ClientConfig struct {
	Timeout            time.Duration
	APIToken           string
	HTTPClient         *http.Client
	MaxRetries         int
	BaseURL            string
	ServerEdition      string
	TrustStore         string
	TrustStorePassword string
	Debug              bool
	// Logger is the logger to send logging messages to.
	Logger logging.LeveledLoggerInterface
}

func NewClient(clientConfig *ClientConfig) (*Client, error) {
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
	if clientConfig.ServerEdition == "" {
		clientConfig.ServerEdition = "community"
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

// ProjectKey returns the SonarQube project key for given context and artifact prefix.
// Monorepo support: separate projects in SonarQube.
// See https://community.sonarsource.com/t/monorepo-and-sonarqube/37990/3.
func ProjectKey(ctxt *pipelinectxt.ODSContext, artifactPrefix string) string {
	sonarProject := fmt.Sprintf("%s-%s", ctxt.Project, ctxt.Component)
	if len(artifactPrefix) > 0 {
		sonarProject = fmt.Sprintf("%s-%s", sonarProject, strings.TrimSuffix(artifactPrefix, "-"))
	}
	return sonarProject
}

func (c *Client) logger() logging.LeveledLoggerInterface {
	return c.clientConfig.Logger
}

func (c *Client) javaSystemProperties() []string {
	return []string{
		fmt.Sprintf("-Djavax.net.ssl.trustStore=%s", c.clientConfig.TrustStore),
		fmt.Sprintf("-Djavax.net.ssl.trustStorePassword=%s", c.clientConfig.TrustStorePassword),
	}
}

func (c *Client) get(urlPath string) (int, []byte, error) {
	u, err := c.baseURL.Parse(urlPath)
	if err != nil {
		return 0, nil, fmt.Errorf("parse URL path: %w", err)
	}
	c.logger().Debugf("GET %s", u)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, nil, fmt.Errorf("could not create request: %s", err)
	}

	res, err := c.do(req)
	if err != nil {
		return 500, nil, fmt.Errorf("got error %s", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	return res.StatusCode, body, err
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	// The user token is sent via the login field of HTTP basic authentication,
	// without any password. See https://docs.sonarqube.org/latest/extend/web-api/.
	credentials := fmt.Sprintf("%s:", c.clientConfig.APIToken)
	basicAuth := b64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth))
	return c.httpClient.Do(req)
}
