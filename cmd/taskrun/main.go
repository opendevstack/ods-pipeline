package main

import (
	"bytes"
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

	namespace := tekton.PrepareConditionsForTaskRun(clientset, storageCapacity, sourceDir, storageClassName, persistentVolumeClaimName)

	applyYAMLFile(namespace, *tektonFilesDir, *taskFileName)
	applyYAMLFile(namespace, *tektonFilesDir, *taskRunFileName)

	// Wait for task to run
	// Show logs
	// Check if task was successful
	// Check local folder
}

func applyYAMLFile(namespace string, fileDir string, fileName string) {

	filePath := fmt.Sprintf("%s/%s", fileDir, fileName)
	stdout, stderr, err := runCmd("kubectl", []string{"-n", namespace, "apply", "-f", filePath})

	fmt.Println(string(stdout))
	fmt.Println(string(stderr))
	check(err)
}

func runCmd(executable string, args []string) (outBytes, errBytes []byte, err error) {
	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outBytes = stdout.Bytes()
	errBytes = stderr.Bytes()
	return outBytes, errBytes, err
}