package main

import (
	"flag"
	"fmt"
	"os/exec"
	// k "github.com/opendevstack/pipeline/pkg/kubernetes"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	taskRunName := flag.String("taskrun-name", "", "Name of the Tekton TaskRun to run")
	sourceDir := flag.String("source-dir", "", "Source directory whose files will be available to the container.")
	outputDir := flag.String("output-dir", "", "Output directory whose container artifacts will be available to inspect in the host.")
	tasksDir := flag.String("tasks-dir", "scripts", "Tasks directory that holds the Tekton task definitions.")

	flag.Parse()

	fmt.Println("taskRunName:", *taskRunName)
	fmt.Println("sourceDir:", *sourceDir)
	fmt.Println("outputDir:", *outputDir)
	fmt.Println("tasksDir:", *tasksDir)

	// clientset := k.NewClient()

	// namespace := "default"
	// volumeName := "test-volume"
	// hostPath := "/files" // check scripts/kind-with-registry.sh
	// storageClassName := "standard"
	// pvcName := "task-pv-claim"
	// capacity := "1Gi"

	// _, err := tekton.CreatePV(clientset, volumeName, capacity, hostPath, storageClassName)
	// check(err)

	// _, err = tekton.CreatePVC(clientset, pvcName, capacity, storageClassName, namespace)
	// check(err)

	// _, err = tekton.StartPodWithPVC(clientset, "task-pv-claim", namespace)
	// check(err)

	// Apply Task
	taskpath := fmt.Sprintf("../../%s/%s", *tasksDir, "task.yaml")
	fmt.Println(taskpath)
	output, err := exec.Command("kubectl", "apply", "-f", taskpath).Output()
	check(err)
	fmt.Println(string(output))

	// Apply TaskRun
	taskrunpath := fmt.Sprintf("../../%s/%s", *tasksDir, "taskrun.yaml")
	fmt.Println(taskrunpath)
	output, err = exec.Command("kubectl", "apply", "-f", taskrunpath).Output()
	check(err)
	fmt.Println(string(output))

	// err = tekton.ExecTaskRun(clientset, *taskRunName)
	// check(err)

}
