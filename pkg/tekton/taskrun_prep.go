package tekton

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreatePV(clientset *kubernetes.Clientset, pvName string, capacity string, hostPath string, storageClassName string) (*v1.PersistentVolume, error) {

	pv, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(),
		&v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name: pvName,
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

func CreatePVC(clientset *kubernetes.Clientset, pvcName string, capacity string, storageClassName string, namespace string) (*v1.PersistentVolumeClaim, error) {

	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(),
		&v1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: pvcName,
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

func StartPodWithPVC(clientset *kubernetes.Clientset, pvcName string, namespace string) (*v1.Pod, error) {

	pod, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "alpine",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: "foo",
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
						}},
				},
			},
			Containers: []v1.Container{
				{
					VolumeMounts:    []v1.VolumeMount{{Name: "foo", MountPath: "/filesincontainer"}},
					Image:           "alpine:latest",
					ImagePullPolicy: v1.PullIfNotPresent,
					Command:         []string{"/bin/sh", "-c", "--"},
					Args:            []string{"while true; do sleep 30; done;"},
					Name:            "alpine"},
			},
		}},
		metav1.CreateOptions{})

	return pod, err
}

func UploadFilesToPod(clientset *kubernetes.Clientset, podName string) {

}

func MountPVCInTaskRun() {

}

func ExecTaskRun(taskRunName string) {

}

func DownloadWorkspaceArtifacts(podName string) {

}
