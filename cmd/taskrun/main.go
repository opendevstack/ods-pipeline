package main

import (
	k "github.com/opendevstack/pipeline/pkg/kubernetes"
	"github.com/opendevstack/pipeline/pkg/tekton"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// nexusURL := flag.String("nexus-url", "", "URL of Nexus instance")
	// nexusUser := flag.String("nexus-user", "", "User of Nexus instance")
	// nexusPassword := flag.String("nexus-password", "", "Password of Nexus user")
	// repository := flag.String("repository", "", "Nexus repository")
	// group := flag.String("group", "", "Repository group")
	// file := flag.String("file", "", "Filename to upload (absolute")

	// flag.Parse()

	clientset := k.NewClient()

	namespace := "default"
	volumeName := "test-volume"
	hostPath := "/files"
	storageClassName := "standard"
	pvcName := "task-pv-claim"

	_, err := tekton.CreatePV(clientset, volumeName, "1Gi", hostPath, storageClassName)
	check(err)

	_, err = tekton.CreatePVC(clientset, pvcName, "1Gi", storageClassName, namespace)
	check(err)

	_, err = tekton.StartPodWithPVC(clientset, "task-pv-claim", namespace)
	check(err)

	// fmt.Printf("Pod Spec: %s\n", &pod.Spec.String())

	// tekton.UploadFilesToPod(clientset, pod.ObjectMeta.Name)
}
