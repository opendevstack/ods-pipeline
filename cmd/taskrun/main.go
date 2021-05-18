package main

import (
	"flag"
	"fmt"
	"os/exec"

	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/tekton"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/* Usage:

go run main.go \
--task-file-name task-hello-world.yaml  \
--taskrun-file-name taskrun-hello-world.yaml \
--pvc-claim-name task-pv-claim \
--source-dir /files

*/
func main() {
	taskFileName := flag.String("task-file-name", "", "Name of the YAML file that defines the Tekton Task.")
	taskRunFileName := flag.String("taskrun-file-name", "", "Name of the YAML file that defines the Tekton TaskRun to run.")
	tektonFilesDir := flag.String("tekton-files-dir", "../../scripts", "Directory where the Tekton YAML files are.")
	sourceDir := flag.String("source-dir", "/files", "Source directory whose files will be available to the container.") // check scripts/kind-with-registry.sh
	storageClassName := flag.String("storage-class-name", "standard", "Storage Class name of the PV and PVC to create.")
	storageCapacity := flag.String("storage-capacity", "1Gi", "Storage capacity of the PV and PVC to create.")
	persistentVolumeClaimName := flag.String("pvc-claim-name", "", "Name of the Persistent Volume Claim defined in the TaskRun.")

	flag.Parse()

	clientset := k.NewClient()

	namespace := tekton.CreateRandomNamespace(clientset)

	_, err := tekton.CreatePersistentVolume(clientset, *storageCapacity, *sourceDir, *storageClassName)
	check(err)

	_, err = tekton.CreatePersistentVolumeClaim(clientset, *persistentVolumeClaimName, *storageCapacity, *storageClassName, namespace)
	check(err)

	taskFilePath := fmt.Sprintf("%s/%s", *tektonFilesDir, *taskFileName)
	applyYAMLFile(namespace, taskFilePath)

	taskRunFilePath := fmt.Sprintf("%s/%s", *tektonFilesDir, *taskRunFileName)
	applyYAMLFile(namespace, taskRunFilePath)
}

func applyYAMLFile(namespace string, filePath string) {

	output, err := exec.Command("kubectl", "-n", namespace, "apply", "-f", filePath).Output()
	check(err)
	fmt.Println(string(output))
}
