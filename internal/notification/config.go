package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/opendevstack/ods-pipeline/internal/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configMapName           = "ods-notification"
	urlProperty             = "url"
	methodProperty          = "method"
	contentTypeProperty     = "contentType"
	requestTemplateProperty = "requestTemplate"
	notifyOnStatusProperty  = "notifyOnStatus"
	enabledProperty         = "enabled"
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
	cm, err := kubernetesClient.GetConfigMap(ctxt, configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load %s ConfigMap: %v", configMapName, err)
	}

	enabledPropValue, ok := cm.Data[enabledProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", configMapName, enabledProperty)
	}

	enabled, err := strconv.ParseBool(enabledPropValue)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %s to bool", enabledPropValue)
	}

	url, ok := cm.Data[urlProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", configMapName, urlProperty)
	}

	method, ok := cm.Data[methodProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", configMapName, methodProperty)
	}

	contentType, ok := cm.Data[contentTypeProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", configMapName, contentTypeProperty)
	}

	notifyOnStatus, ok := cm.Data[notifyOnStatusProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specifiy '%s' property", configMapName, notifyOnStatusProperty)
	}

	decoder := json.NewDecoder(strings.NewReader(notifyOnStatus))
	var notificationStatusValues []string
	err = decoder.Decode(&notificationStatusValues)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("decoding notification status properties failed: %w", err)
	}

	text, ok := cm.Data[requestTemplateProperty]
	if !ok {
		return nil, fmt.Errorf("%s doesn't specify '%s' property", configMapName, requestTemplateProperty)
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
