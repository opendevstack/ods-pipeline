package manager

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Annotation to set the storage provisioner for a PVC.
	storageProvisionerAnnotation = "volume.beta.kubernetes.io/storage-provisioner"
	// PVC finalizer.
	pvcProtectionFinalizer = "kubernetes.io/pvc-protection"
)

// createPVCIfRequired if it does not exist yet
func (s *Server) createPVCIfRequired(ctxt context.Context, pData PipelineData) error {
	_, err := s.KubernetesClient.GetPersistentVolumeClaim(ctxt, pData.PVC, metav1.GetOptions{})
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return fmt.Errorf("could not determine if %s already exists: %w", pData.PVC, err)
		}
		vm := corev1.PersistentVolumeFilesystem
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:        pData.PVC,
				Labels:      map[string]string{repositoryLabel: pData.Repository},
				Finalizers:  []string{pvcProtectionFinalizer},
				Annotations: map[string]string{},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(s.StorageConfig.Size),
					},
				},
				StorageClassName: &s.StorageConfig.ClassName,
				VolumeMode:       &vm,
			},
		}
		if s.StorageConfig.Provisioner != "" {
			pvc.Annotations[storageProvisionerAnnotation] = s.StorageConfig.Provisioner
		}
		_, err := s.KubernetesClient.CreatePersistentVolumeClaim(ctxt, pvc, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func makePVCName(component string) string {
	pvcName := fmt.Sprintf("ods-workspace-%s", strings.ToLower(component))
	return fitStringToMaxLength(pvcName, 63) // K8s label max length to be on the safe side.
}
