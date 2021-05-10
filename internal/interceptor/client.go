package interceptor

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	namespaceFile     = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	tokenFile         = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	caCert            = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	pipelineFilename  = "ods.yml"
	tektonAPIBasePath = "/apis/tekton.dev/v1beta1"
	letterBytes       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	namespaceSuffix   = "-cd"
	apiHostEnvVar     = "API_HOST"
	apiHostDefault    = "openshift.default.svc.cluster.local"
	repoBaseEnvVar    = "REPO_BASE"
	tokenEnvVar       = "ACCESS_TOKEN"
)

// Client makes requests, e.g. to create and delete pipelines, or to forward
// event payloads.
type Client interface {
	GetPipelineResourceVersion(name string) (int, error)
	ApplyPipeline(pipelineBody []byte, data PipelineData) (int, error)
}

type ocClient struct {
	HTTPClient *http.Client
	APIBaseURL string
	Namespace  string
	Token      string
}

type buildConfig struct {
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
}

func (c *ocClient) ApplyPipeline(pipelineBody []byte, data PipelineData) (int, error) {
	url := fmt.Sprintf(
		"%s/namespaces/%s/pipelines",
		c.APIBaseURL,
		c.Namespace,
	)
	verb := "POST"
	if data.ResourceVersion > 0 {
		verb = "PUT"
		url = fmt.Sprintf("%s/%s", url, data.Name)
	}
	req, _ := http.NewRequest(verb, url, bytes.NewBuffer(pipelineBody))
	res, err := c.do(req)
	if err != nil {
		return 500, fmt.Errorf("could not make OpenShift request: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 500, fmt.Errorf("could not read OpenShift response body: %s", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return res.StatusCode, fmt.Errorf("unexpected return code for %s %s: %s", verb, url, body)
	}

	return res.StatusCode, nil
}

// GetPipeline determines whether the pipeline corresponding to the given
// event already exists.
func (c *ocClient) GetPipelineResourceVersion(name string) (int, error) {
	resourceVersion := 0
	url := fmt.Sprintf(
		"%s/namespaces/%s/pipelines/%s",
		c.APIBaseURL,
		c.Namespace,
		name,
	)
	req, _ := http.NewRequest(
		"GET",
		url,
		nil,
	)
	res, err := c.do(req)
	if err != nil {
		return resourceVersion, fmt.Errorf("could not make OpenShift request: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return resourceVersion, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resourceVersion, fmt.Errorf("could not read OpenShift response: %s", err)
	}

	if res.StatusCode != 200 {
		return resourceVersion, fmt.Errorf("unexpected OpenShift response: [%d] %s", res.StatusCode, string(body))
	}

	bc := buildConfig{}
	err = json.Unmarshal(body, &bc)
	if err != nil {
		return resourceVersion, fmt.Errorf("could not parse existing pipeline: %s", err)
	}
	rv, err := strconv.Atoi(bc.Metadata.ResourceVersion)
	if err != nil {
		return resourceVersion, fmt.Errorf("resourceVersion is not an int: %s", err)
	}
	return rv, nil
}

// do executes the request.
func (c *ocClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	return c.HTTPClient.Do(req)
}

// NewServer returns a new OpenShift client.
func NewClient(APIHost string, namespace string) (*ocClient, error) {
	token, err := getFileContent(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("could not get token: %s", err)
	}

	secureClient, err := getSecureClient()
	if err != nil {
		return nil, fmt.Errorf("could not get client: %s", err)
	}

	baseURL := fmt.Sprintf(
		"https://%s%s",
		APIHost,
		tektonAPIBasePath,
	)

	return &ocClient{
		HTTPClient: secureClient,
		APIBaseURL: baseURL,
		Token:      token,
		Namespace:  namespace,
	}, nil
}

func getSecureClient() (*http.Client, error) {
	// Load CA cert
	caCert, err := ioutil.ReadFile(caCert)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport, Timeout: 10 * time.Second}, nil
}

func getFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
