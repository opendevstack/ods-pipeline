package test

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/tekton"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestSatisfiedBy(t *testing.T) {

	clients := k.NewClients()
	tektonTasksDir := "../scripts" // Directory where the Tekton YAML files are
	taskFileName := "task-hello-world.yaml"
	taskName := "hello-world"
	storageCapacity := "1Gi"
	storageClassName := "" // if using KinD, set it to "standard"

	tests := map[string]struct {
		sourceDirectory string
		workspaceName   string
		claimName       string
		params          map[string]string
		wantSuccess     string
	}{
		"task should return success": {
			sourceDirectory: "/mnt/c/src/ods-pipeline/test/goapp", //  for a local volume name, only \"[a-zA-Z0-9][a-zA-Z0-9_.-]\" are allowed. If you inte ││ nded to pass a host directory, use absolute path
			workspaceName:   "source",                             // must exist in the Task definition
			claimName:       "task-pv-claim",
			params:          map[string]string{"message": "foo"},
			wantSuccess:     "True",
		},
	}

	for name, tc := range tests {

		// setup code
		// It is assumed that Tekton is already installed in the KinD cluster.
		// If not, run scripts/kind-with-registry.sh
		// Prior to run the test for the Task, we create:
		// - A local temporary directory.
		// - A Persistent Volume (PV) with hostPath pointing to the local temp dir.
		// - A Persistent Volume Claim (PVC) that will be referenced in the TaskRun to mount the local temp dir.
		namespace := tekton.PrepareConditionsForTaskRun(clients.KubernetesClientSet, &storageCapacity, &tc.sourceDirectory, &storageClassName, &tc.claimName)
		applyYAMLFile(namespace, tektonTasksDir, taskFileName)

		t.Run(name, func(t *testing.T) {
			actual, err := tekton.Run(clients.TektonClientSet, taskName, tc.params, tc.workspaceName, tc.claimName, namespace)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Printf("Status: %s\n", actual.Status.Status.Conditions[0].Status)
			fmt.Printf("Reason: %s\n", actual.Status.Status.Conditions[0].Reason)

			status := string(actual.Status.Status.Conditions[0].Status)

			if status != tc.wantSuccess {
				t.Errorf("Got: %+v, want: %+v.", status, tc.wantSuccess)
			}
		})

		// TODO: tear-down code
	}
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
