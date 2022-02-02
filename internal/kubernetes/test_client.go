package kubernetes

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kschema "k8s.io/apimachinery/pkg/runtime/schema"
)

// TestClient returns mocked resources.
type TestClient struct {
	// PersistentVolumeClaims is the pool of pipelines which can be retrieved.
	PVCs []*corev1.PersistentVolumeClaim
	// FailCreatePVC lets PVC creation fail.
	FailCreatePVC bool
	// CreatedPVCs is a slice of created PVC names.
	CreatedPVCs []string
	// ConfigMaps which can be retrieved
	CMs []*corev1.ConfigMap
}

func (c *TestClient) GetPersistentVolumeClaim(ctxt context.Context, name string, options metav1.GetOptions) (*corev1.PersistentVolumeClaim, error) {
	for _, p := range c.PVCs {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("pipeline %s not found", name)
}

func (c *TestClient) CreatePersistentVolumeClaim(ctxt context.Context, pipeline *corev1.PersistentVolumeClaim, options metav1.CreateOptions) (*corev1.PersistentVolumeClaim, error) {
	c.CreatedPVCs = append(c.CreatedPVCs, pipeline.Name)
	if c.FailCreatePVC {
		return nil, errors.New("creation error")
	}
	return pipeline, nil
}

func (c *TestClient) GetConfigMap(ctxt context.Context, cmName string, options metav1.GetOptions) (*corev1.ConfigMap, error) {
	for _, cm := range c.CMs {
		if cm.Name == cmName {
			return cm, nil
		}
	}
	return nil, kerrors.NewNotFound(kschema.GroupResource{
		Group:    "core",
		Resource: "ConfigMap",
	}, cmName)
}

func (c *TestClient) GetConfigMapKey(ctxt context.Context, cmName, key string, options metav1.GetOptions) (string, error) {
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
