package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	Enabled        bool
	URL            string
	Method         string
	ContentType    string
	NotifyOnStatus []string
	Template       *template.Template
}

func ReadConfigFromConfigMap(ctxt context.Context, kubernetesClient kubernetes.ClientInterface) (*Config, error) {
	cm, err := kubernetesClient.GetConfigMap(ctxt, ConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load %s ConfigMap: %v", ConfigMapName, err)
	}

	enabledPropValue, ok := cm.Data[EnabledProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", ConfigMapName, EnabledProperty)
	}

	enabled, err := strconv.ParseBool(enabledPropValue)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %s to bool", enabledPropValue)
	}

	url, ok := cm.Data[UrlProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", ConfigMapName, UrlProperty)
	}

	method, ok := cm.Data[MethodProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", ConfigMapName, MethodProperty)
	}

	contentType, ok := cm.Data[ContentTypeProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", ConfigMapName, ContentTypeProperty)
	}

	notifyOnStatus, ok := cm.Data[NotifyOnStatusProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specifiy '%s' property", ConfigMapName, NotifyOnStatusProperty)
	}

	decoder := json.NewDecoder(strings.NewReader(notifyOnStatus))
	var notificationStatusValues []string
	err = decoder.Decode(&notificationStatusValues)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("decoding notification status properties failed: %w", err)
	}

	text, ok := cm.Data[RequestTemplateProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", ConfigMapName, RequestTemplateProperty)
	}

	requestTemplate, err := template.New("requestTemplate").Parse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse requestTemplate Template")
	}

	return &Config{
		enabled,
		url,
		method,
		contentType,
		notificationStatusValues,
		requestTemplate,
	}, nil
}
