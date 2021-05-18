package main

import (
	"flag"

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

	flag.Parse()

	clientset := k.NewClient()

	namespace := "default"
	_, err := tekton.CreatePVC(clientset, "task-pv-claim", namespace)
	check(err)

	pod, err := tekton.StartPodWithPVC(clientset, "task-pv-claim", namespace)
	check(err)

	tekton.UploadFilesToPod(clientset, pod.ObjectMeta.Name)
}
