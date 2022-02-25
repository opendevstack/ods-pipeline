package kubernetes

import (
	"errors"
	"path/filepath"

	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Clients struct {
	KubernetesClientSet *kubernetes.Clientset
	TektonClientSet     *tekton.Clientset
}

func NewClients() *Clients {
	// TODO: make configurable from outside
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the Kubernetes clientset
	kubernetesClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create the Tekton clientset
	tektonClientSet, err := tekton.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Clients{
		KubernetesClientSet: kubernetesClientset,
		TektonClientSet:     tektonClientSet,
	}
}

func NewInClusterClientset() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	return kubernetes.NewForConfig(config)
}

// Client represents a Kubernetes client, wrapping
// k8s.io/client-go/kubernetes.Clientset
type Client struct {
	clientConfig *ClientConfig
}

// ClientConfig configures a Tekton client.
type ClientConfig struct {
	// Kubernetes namespace.
	Namespace string
	// Logger is the logger to send logging messages to.
	Logger logging.LeveledLoggerInterface
	// KubernetesClient is the wrapped Kubernetes client.
	KubernetesClient *kubernetes.Clientset
}

type ClientInterface interface {
	ClientPersistentVolumeClaimInterface
	ClientConfigMapInterface
}

// NewInClusterClient initializes a Kubernetes client from within a cluster.
func NewInClusterClient(clientConfig *ClientConfig) (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// create the Kubernetes clientset
	kubernetesClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	clientConfig.KubernetesClient = kubernetesClientSet
	return NewClient(clientConfig)
}

// NewClient initializes a Kubernetes client.
func NewClient(clientConfig *ClientConfig) (*Client, error) {
	if clientConfig.Namespace == "" {
		return nil, errors.New("namespace is required")
	}

	if clientConfig.KubernetesClient == nil {
		return nil, errors.New("kubernetes client is required")
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

func (c *Client) coreV1Client() clientCoreV1.CoreV1Interface {
	return c.clientConfig.KubernetesClient.CoreV1()
}

func (c *Client) persistentVolumeClaimsClient() clientCoreV1.PersistentVolumeClaimInterface {
	return c.coreV1Client().PersistentVolumeClaims(c.namespace())
}

func (c *Client) configMapsClient() clientCoreV1.ConfigMapInterface {
	return c.coreV1Client().ConfigMaps(c.namespace())
}
