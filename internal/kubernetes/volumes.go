package kubernetes

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreatePersistentVolume(clientset kubernetes.Interface, pvName string, capacity string, hostPath string, storageClassName string) (*v1.PersistentVolume, error) {

	log.Printf("Create persistent volume %s", pvName)

	pv, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(),
		&v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:   pvName,
				Labels: map[string]string{"app.kubernetes.io/managed-by": "ods-pipeline"},
			},
			Spec: v1.PersistentVolumeSpec{
				Capacity: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse(capacity),
				},
				AccessModes:                   []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				PersistentVolumeSource:        v1.PersistentVolumeSource{HostPath: &v1.HostPathVolumeSource{Path: hostPath}},
				PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimRetain,
				StorageClassName:              storageClassName,
			},
		}, metav1.CreateOptions{})

	return pv, err
}

func CreatePersistentVolumeClaim(clientset kubernetes.Interface, capacity string, storageClassName string, namespace string) (*v1.PersistentVolumeClaim, error) {

	pvcName := "task-pv-claim"
	log.Printf("Create persistent volume claim %s", pvcName)

	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(),
		&v1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:   pvcName,
				Labels: map[string]string{"app.kubernetes.io/managed-by": "ods-pipeline"},
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				StorageClassName: &storageClassName,
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceName(v1.ResourceStorage): resource.MustParse(capacity),
					},
				},
			},
		}, metav1.CreateOptions{})

	return pvc, err
}
