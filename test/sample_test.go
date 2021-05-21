package test

import (
	"fmt"
	"testing"

	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/tekton"
)

func TestSatisfiedBy(t *testing.T) {

	clients := k.NewClients()
	taskName := "hello-world"
	storageCapacity := "1Gi"
	storageClassName := "" //"standard"

	tests := map[string]struct {
		sourceDirectory string
		workspaceName   string
		claimName       string
		params          map[string]string
		wantSuccess     string
	}{
		"task should return success": {
			sourceDirectory: "./go-app",
			workspaceName:   "source", // must exist in the Task definition
			claimName:       "task-pv-claim",
			params:          map[string]string{"message": "foo"},
			wantSuccess:     "True",
		},
	}

	for name, tc := range tests {

		// setup code
		namespace := tekton.PrepareConditionsForTaskRun(clients.KubernetesClientSet, &storageCapacity, &tc.sourceDirectory, &storageClassName, &tc.claimName)

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
