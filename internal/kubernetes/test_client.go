package kubernetes

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestClient returns mocked resources.
type TestClient struct {
	// PersistentVolumeClaims is the pool of pipelines which can be retrieved.
	PVCs []*corev1.PersistentVolumeClaim
	// FailCreatePVC lets PVC creation fail.
	FailCreatePVC bool
	// CreatedPVCs is a slice of created PVC names.
	CreatedPVCs []string
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
