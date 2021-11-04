package installation

import (
	"context"
	"fmt"

	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "k8s.io/client-go/kubernetes"
)

const (
	NexusConfigMapName                   = "ods-nexus"
	NexusSecretName                      = "ods-nexus-auth"
	NexusSecretUsernameKey               = "username"
	NexusSecretPasswordKey               = "password"
	NexusConfigMapURLKey                 = "url"
	NexusConfigMapTemporaryRepositoryKey = "temporaryRepository"
	NexusConfigMapPermanentRepositoryKey = "permanentRepository"
)

type NexusRepositories struct {
	Temporary string
	Permanent string
}

// NewNexusClientConfig returns a *nexus.ClientConfig which is derived
// from the information about Nexus located in the given Kubernetes namespace.
func NewNexusClientConfig(c *kclient.Clientset, namespace string, logger logging.LeveledLoggerInterface) (*nexus.ClientConfig, error) {
	nexusSecret, err := c.CoreV1().Secrets(namespace).
		Get(context.TODO(), NexusSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get Nexus secret: %w", err)
	}
	nexusConfigMap, err := c.CoreV1().ConfigMaps(namespace).
		Get(context.TODO(), NexusConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get Nexus config: %w", err)
	}
	nexusClientConfig := &nexus.ClientConfig{
		Username: string(nexusSecret.Data[NexusSecretUsernameKey]),
		Password: string(nexusSecret.Data[NexusSecretPasswordKey]),
		BaseURL:  nexusConfigMap.Data[NexusConfigMapURLKey],
		Logger:   logger,
	}
	return nexusClientConfig, nil
}

func GetNexusRepositories(c *kclient.Clientset, namespace string) (*NexusRepositories, error) {
	nexusConfigMap, err := c.CoreV1().ConfigMaps(namespace).
		Get(context.TODO(), NexusConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get Nexus config: %w", err)
	}
	return &NexusRepositories{
		Temporary: nexusConfigMap.Data[NexusConfigMapTemporaryRepositoryKey],
		Permanent: nexusConfigMap.Data[NexusConfigMapPermanentRepositoryKey],
	}, nil
}
