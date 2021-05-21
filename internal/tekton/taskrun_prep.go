package tekton

import (
	"os"

	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"k8s.io/client-go/kubernetes"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func PrepareConditionsForTaskRun(clientset *kubernetes.Clientset, storageCapacity *string, sourceDir *string, storageClassName *string, persistentVolumeClaimName *string) string {

	// Create temp directory
	if _, err := os.Stat(*sourceDir); os.IsNotExist(err) {
		os.Mkdir(*sourceDir, os.ModeDir)
	}

	namespace := k.CreateRandomNamespace(clientset)

	_, err := k.CreatePersistentVolume(clientset, *storageCapacity, *sourceDir, *storageClassName)
	check(err)

	_, err = k.CreatePersistentVolumeClaim(clientset, *persistentVolumeClaimName, *storageCapacity, *storageClassName, namespace)
	check(err)

	return namespace
}
