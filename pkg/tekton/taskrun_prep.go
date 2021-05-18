package tekton

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreatePVC(clientset *kubernetes.Clientset, pvcName string, namespace string) (*v1.PersistentVolumeClaim, error) {

	var standardClass = "standard"

	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(),
		&v1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: pvcName,
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				StorageClassName: &standardClass,
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
			Volumes: []v1.Volume{{
				Name: "test-volume",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
						ReadOnly:  false,
					},
				},
			}},
			Containers: []v1.Container{
				{
					Image:   "alpine:latest",
					Command: []string{"/bin/sh", "-c", "--"},
					Args:    []string{"while true; do sleep 30; done;"},
					Name:    "alpine"}},
		},
	}, metav1.CreateOptions{
		TypeMeta:     metav1.TypeMeta{},
		DryRun:       []string{},
		FieldManager: "",
	})

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
