package installation

import (
	"context"
	"fmt"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "k8s.io/client-go/kubernetes"
)

const (
	BitbucketConfigMapName     = "ods-bitbucket"
	BitbucketSecretName        = "ods-bitbucket-auth"
	BitbucketSecretAPITokenKey = "password"
	BitbucketConfigMapURLKey   = "url"
)

// NewBitbucketClientConfig returns a *bitbucket.ClientConfig which is derived
// from the information about Bitbucket located in the given Kubernetes namespace.
func NewBitbucketClientConfig(c *kclient.Clientset, namespace string, logger logging.LeveledLoggerInterface) (*bitbucket.ClientConfig, error) {
	bitbucketSecret, err := c.CoreV1().Secrets(namespace).
		Get(context.TODO(), BitbucketSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get Bitbucket secret: %w", err)
	}
	bitbucketConfigMap, err := c.CoreV1().ConfigMaps(namespace).
		Get(context.TODO(), BitbucketConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get Bitbucket config: %w", err)
	}
	bitbucketClientConfig := &bitbucket.ClientConfig{
		APIToken: string(bitbucketSecret.Data[BitbucketSecretAPITokenKey]),
		BaseURL:  bitbucketConfigMap.Data[BitbucketConfigMapURLKey],
		Logger:   logger,
	}
	return bitbucketClientConfig, nil
}
