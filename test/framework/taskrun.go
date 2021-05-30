package framework

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelineclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testdataWorkspacePath = "testdata/workspaces"

func CreateTaskRunWithParams(tknClient *pipelineclientset.Clientset, taskRefKind string, taskName string, parameters map[string]string, workspaces map[string]string, namespace string) (*tekton.TaskRun, error) {

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

	var tk tekton.TaskKind
	switch taskRefKind {
	case string(tekton.ClusterTaskKind):
		tk = tekton.ClusterTaskKind
	case string(tekton.NamespacedTaskKind):
		tk = tekton.NamespacedTaskKind
	default:
		log.Fatalf("Don't know type %s\n", taskRefKind)
	}

	taskWorkspaces := []tekton.WorkspaceBinding{}
	for wn, wd := range workspaces {
		wsDirName := filepath.Base(wd)
		taskWorkspaces = append(taskWorkspaces, tekton.WorkspaceBinding{
			Name: wn,
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "task-pv-claim",
				ReadOnly:  false,
			},
			SubPath: filepath.Join(testdataWorkspacePath, wsDirName),
		})
	}

	tr, err := tknClient.TektonV1beta1().TaskRuns(namespace).Create(context.TODO(),
		&tekton.TaskRun{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-taskrun-%s", taskName, uuid.NewV4()),
			},
			Spec: tekton.TaskRunSpec{
				TaskRef:    &tekton.TaskRef{Kind: tk, Name: taskName},
				Params:     tektonParams,
				Workspaces: taskWorkspaces,
			},
		},
		metav1.CreateOptions{})

	return tr, err
}

func getTr(ctx context.Context, t *testing.T, c pipelineclientset.Interface, name, ns string) (tr *tekton.TaskRun) {
	t.Helper()
	tr, err := c.TektonV1beta1().TaskRuns(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	return tr
}

type conditionFn func(*tekton.TaskRun) bool

func WaitForCondition(ctx context.Context, t *testing.T, c pipelineclientset.Interface, name, ns string, cond conditionFn, timeout time.Duration) *tekton.TaskRun {
	t.Helper()

	// Do a first quick check before setting the watch
	tr := getTr(ctx, t, c, name, ns)
	if cond(tr) {
		return tr
	}

	w, err := c.TektonV1beta1().TaskRuns(ns).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}))
	if err != nil {
		t.Errorf("error watching taskrun: %s", err)
	}

	// Setup a timeout channel
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		timeoutChan <- struct{}{}
	}()

	// Wait for the condition to be true or a timeout
	for {
		select {
		case ev := <-w.ResultChan():
			if ev.Object != nil {
				tr := ev.Object.(*tekton.TaskRun)
				if cond(tr) {
					return tr
				}
			}
		case <-timeoutChan:
			t.Fatal("time out")
		}
	}
}

func Done(tr *tekton.TaskRun) bool {
	return tr.IsDone()
}
