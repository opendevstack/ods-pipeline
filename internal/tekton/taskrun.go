package tekton

import (
	"context"
	"fmt"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	versioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Run(tknClient *versioned.Clientset, taskName string, parameters map[string]string, workspaceName string, claimName string, namespace string) (*tekton.TaskRun, error) {

	var tektonParams []tekton.Param

	for _, p := range tektonParams {

		tektonParams = append(tektonParams, tekton.Param{
			Name:  p.Name,
			Value: tekton.ArrayOrString{StringVal: p.Value.StringVal},
		})

	}

	fmt.Printf("Tekton params: %s\n", tektonParams)

	tr, err := tknClient.TektonV1beta1().TaskRuns(namespace).Create(context.TODO(),
		&tekton.TaskRun{
			Spec: tekton.TaskRunSpec{
				Params: tektonParams,
				Workspaces: []tekton.WorkspaceBinding{
					{
						Name: workspaceName,
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: claimName,
							ReadOnly:  false,
						},
					},
				},
			},
		},
		v1.CreateOptions{})

	return tr, err
}
