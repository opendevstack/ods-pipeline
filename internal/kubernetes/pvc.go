package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientPersistentVolumeClaimInterface interface {
	GetPersistentVolumeClaim(ctxt context.Context, name string, options metav1.GetOptions) (*corev1.PersistentVolumeClaim, error)
	CreatePersistentVolumeClaim(ctxt context.Context, pipeline *corev1.PersistentVolumeClaim, options metav1.CreateOptions) (*corev1.PersistentVolumeClaim, error)
}

func (c *Client) GetPersistentVolumeClaim(ctxt context.Context, name string, options metav1.GetOptions) (*corev1.PersistentVolumeClaim, error) {
	c.logger().Debugf("Get persistent volume claim %s", name)
	return c.persistentVolumeClaimsClient().Get(ctxt, name, options)
}

func (c *Client) CreatePersistentVolumeClaim(ctxt context.Context, pvc *corev1.PersistentVolumeClaim, options metav1.CreateOptions) (*corev1.PersistentVolumeClaim, error) {
	c.logger().Debugf("Create persistent volume claim %s", pvc.Name)
	return c.persistentVolumeClaimsClient().Create(ctxt, pvc, options)
}
