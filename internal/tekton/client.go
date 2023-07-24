package tekton

import (
	"errors"

	"github.com/opendevstack/ods-pipeline/pkg/logging"
	tektonClient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	v1beta1 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1"
	"k8s.io/client-go/rest"
)

// Client represents a Tekton client, wrapping
// github.com/tektoncd/pipeline/pkg/client/clientset/versioned.Clientset
type Client struct {
	clientConfig *ClientConfig
}

// ClientConfig configures a Tekton client.
type ClientConfig struct {
	// Kubernetes namespace.
	Namespace string
	// Logger is the logger to send logging messages to.
	Logger logging.LeveledLoggerInterface
	// TektonClient is the wrapped Tekton client.
	TektonClient *tektonClient.Clientset
}

type ClientInterface interface {
	ClientPipelineRunInterface
}

// NewInClusterClient initializes a Tekton client from within a cluster.
func NewInClusterClient(clientConfig *ClientConfig) (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// create the Tekton clientset
	tektonClientSet, err := tektonClient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	clientConfig.TektonClient = tektonClientSet
	return NewClient(clientConfig)
}

// NewClient initializes a Tekton client.
func NewClient(clientConfig *ClientConfig) (*Client, error) {
	if clientConfig.Namespace == "" {
		return nil, errors.New("namespace is required")
	}

	if clientConfig.TektonClient == nil {
		return nil, errors.New("tekton client is required")
	}

	// Be careful not to pass a variable of type *logging.LeveledLogger
	// holding a nil value. If you pass nil for Logger, make sure it is of
	// logging.LeveledLoggerInterface type.
	if clientConfig.Logger == nil {
		clientConfig.Logger = &logging.LeveledLogger{Level: logging.LevelError}
	}

	return &Client{clientConfig: clientConfig}, nil
}

func (c *Client) logger() logging.LeveledLoggerInterface {
	return c.clientConfig.Logger
}

func (c *Client) namespace() string {
	return c.clientConfig.Namespace
}

func (c *Client) tektonV1beta1Client() v1beta1.TektonV1beta1Interface {
	return c.clientConfig.TektonClient.TektonV1beta1()
}

func (c *Client) pipelineRunsClient() v1beta1.PipelineRunInterface {
	return c.tektonV1beta1Client().PipelineRuns(c.namespace())
}
