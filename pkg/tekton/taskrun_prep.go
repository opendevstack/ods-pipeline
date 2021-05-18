package tekton

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateRandomNamespace(clientset *kubernetes.Clientset) string {

	id := uuid.NewV4()

	ns, err := clientset.CoreV1().Namespaces().Create(context.TODO(),
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: id.String(),
			},
		},
		metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created Namespace: %s\n", ns.Name)

	return ns.Name
}

func CreatePersistentVolume(clientset *kubernetes.Clientset, capacity string, hostPath string, storageClassName string) (*v1.PersistentVolume, error) {

	pvName := uuid.NewV4().String()

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

	fmt.Printf("Created Persistent Volume: %s\n", pv.Name)

	return pv, err
}

func CreatePersistentVolumeClaim(clientset *kubernetes.Clientset, pvcName string, capacity string, storageClassName string, namespace string) (*v1.PersistentVolumeClaim, error) {

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
