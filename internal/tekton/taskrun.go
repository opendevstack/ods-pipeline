package tekton

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	versioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Run(tknClient *versioned.Clientset, taskName string, parameters map[string]string, workspaceName string, claimName string, namespace string) (*tekton.TaskRun, error) {

	var tektonParams []tekton.Param

	for key, value := range parameters {

		tektonParams = append(tektonParams, tekton.Param{
			Name: key,
			Value: tekton.ArrayOrString{
				Type:      "string", // we only provide support to string params for now
				StringVal: value,
			},
		})

	}

	tr, err := tknClient.TektonV1beta1().TaskRuns(namespace).Create(context.TODO(),
		&tekton.TaskRun{
			ObjectMeta: v1.ObjectMeta{
				Name: fmt.Sprintf("%s-taskrun-%s", taskName, uuid.NewV4()),
			},
			Spec: tekton.TaskRunSpec{
				TaskRef: &tekton.TaskRef{Kind: "Task", Name: taskName},
				Params:  tektonParams,
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
