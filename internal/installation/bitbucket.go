package installation

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
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
func NewBitbucketClientConfig(c *kclient.Clientset, namespace string, logger logging.LeveledLoggerInterface, privateCert string) (*bitbucket.ClientConfig, error) {
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
	httpClient := &http.Client{}
	if privateCert != "" {
		cas, err := rootCAs(privateCert)
		if err != nil {
			return nil, err
		}
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: cas},
		}
	}
	bitbucketClientConfig := &bitbucket.ClientConfig{
		APIToken:   string(bitbucketSecret.Data[BitbucketSecretAPITokenKey]),
		BaseURL:    bitbucketConfigMap.Data[BitbucketConfigMapURLKey],
		Logger:     logger,
		HTTPClient: httpClient,
	}
	return bitbucketClientConfig, nil
}

// rootCAs adds localCertFile to the system trusted certs.
func rootCAs(localCertFile string) (*x509.CertPool, error) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := os.ReadFile(localCertFile)
	if err != nil {
		return nil, fmt.Errorf("append %q to RootCAs: %v", localCertFile, err)
	}
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		return nil, fmt.Errorf("no certs appended, using system certs only")
	}
	return rootCAs, nil
}
