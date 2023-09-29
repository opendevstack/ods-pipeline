package installation

import (
	"context"
	"fmt"

	"github.com/opendevstack/ods-pipeline/pkg/logging"
	"github.com/opendevstack/ods-pipeline/pkg/nexus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "k8s.io/client-go/kubernetes"
)

const (
	NexusConfigMapName     = "ods-nexus"
	NexusSecretName        = "ods-nexus-auth"
	NexusSecretUsernameKey = "username"
	NexusSecretPasswordKey = "password"
	NexusConfigMapURLKey   = "url"
)

// NewNexusClientConfig returns a *nexus.ClientConfig which is derived
// from the information about Nexus located in the given Kubernetes namespace.
func NewNexusClientConfig(c kclient.Interface, namespace string, logger logging.LeveledLoggerInterface) (*nexus.ClientConfig, error) {
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
