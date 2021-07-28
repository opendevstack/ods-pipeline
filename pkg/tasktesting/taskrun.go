package tasktesting

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/random"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelineclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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
				Name: fmt.Sprintf("%s-taskrun-%s", taskName, random.PseudoString()),
			},
			Spec: tekton.TaskRunSpec{
				TaskRef:            &tekton.TaskRef{Kind: tk, Name: taskName},
				Params:             tektonParams,
				Workspaces:         taskWorkspaces,
				ServiceAccountName: "pipeline",
			},
		},
		metav1.CreateOptions{})

	return tr, err
}

func waitForTaskRunDone(
	ctx context.Context,
	t *testing.T,
	c pipelineclientset.Interface,
	name, ns string,
	errs chan error,
	done chan *tekton.TaskRun) {

	deadline, _ := ctx.Deadline()
	timeout := time.Until(deadline)
	log.Printf("Waiting up to %v seconds for task %s in namespace %s to be done...\n", timeout.Round(time.Second).Seconds(), name, ns)

	t.Helper()

	w, err := c.TektonV1beta1().TaskRuns(ns).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}))
	if err != nil {
		errs <- fmt.Errorf("error watching taskrun: %s", err)
		return
	}

	// Wait for the TaskRun to be done
	for {
		ev := <-w.ResultChan()
		if ev.Object != nil {
			tr, ok := ev.Object.(*tekton.TaskRun)
			if ok {
				if tr.IsDone() {
					done <- tr
					close(done)
					return
				}
			}

		}
	}
}

func waitForTaskRunPod(
	ctx context.Context,
	c *kubernetes.Clientset,
	taskRunName,
	namespace string,
	podAdded chan *v1.Pod) {
	log.Printf("Waiting for pod related to TaskRun %s to be added to the cluster\n", taskRunName)
	stop := make(chan struct{})

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(c, time.Second*30)
	podsInformer := kubeInformerFactory.Core().V1().Pods().Informer()

	var taskRunPod *v1.Pod

	podsInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				// when a new task is created, watch its events
				pod := obj.(*v1.Pod)
				if strings.HasPrefix(pod.Name, taskRunName) {
					taskRunPod = pod
					log.Printf("TaskRun %s added pod %s to the cluster", taskRunName, pod.Name)
					stop <- struct{}{}
				}

			},
		})

	defer close(stop)
	kubeInformerFactory.Start(stop)

	<-stop
	podAdded <- taskRunPod
}
