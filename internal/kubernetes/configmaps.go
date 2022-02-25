package kubernetes

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientConfigMapInterface interface {
	GetConfigMap(ctxt context.Context, cmName string, options metav1.GetOptions) (*v1.ConfigMap, error)
	GetConfigMapKey(ctxt context.Context, cmName, key string, options metav1.GetOptions) (string, error)
}

func (c *Client) GetConfigMap(ctxt context.Context, cmName string, options metav1.GetOptions) (*v1.ConfigMap, error) {
	c.logger().Debugf("Get configmap %s", cmName)

	return c.configMapsClient().Get(ctxt, cmName, options)
}

func (c *Client) GetConfigMapKey(ctxt context.Context, cmName, key string, options metav1.GetOptions) (string, error) {
	c.logger().Debugf("Get configmap %s", cmName)

	cm, err := c.GetConfigMap(ctxt, cmName, options)
	if err != nil {
		return "", err
	}

	v, ok := cm.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	return v, err
}
