package tasktesting

import (
	"context"
	"errors"
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

func getTr(ctx context.Context, t *testing.T, c pipelineclientset.Interface, name, ns string) (tr *tekton.TaskRun) {
	t.Helper()
	tr, err := c.TektonV1beta1().TaskRuns(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	return tr
}

type conditionFn func(*tekton.TaskRun) bool

func waitForTaskRunDone(
	ctx context.Context,
	t *testing.T,
	c pipelineclientset.Interface,
	name, ns string,
	timeout time.Duration,
	errs chan error,
	done chan bool) {

	log.Printf("Waiting up to %v seconds for task %s in namespace %s to be done...\n", timeout.Seconds(), name, ns)

	t.Helper()

	// Do a first quick check before setting the watch
	// tr := getTr(ctx, t, c, name, ns)
	// if tr.IsDone() {
	// 	return tr, nil
	// }

	w, err := c.TektonV1beta1().TaskRuns(ns).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}))
	if err != nil {
		errs <- fmt.Errorf("error watching taskrun: %s", err)
		return
	}

	// Setup a timeout channel
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		timeoutChan <- struct{}{}
	}()

	// Wait for the TaskRun to be done or time out,
	// or a failure in the pod's events,
	// or the pod's containers to be ready
	for {
		select {
		case ev := <-w.ResultChan():
			if ev.Object != nil {
				tr := ev.Object.(*tekton.TaskRun)
				if tr.IsDone() {
					done <- true
					close(done)
					return
				}
			}

		case err := <-errs:
			if err != nil {
				errs <- fmt.Errorf("Stopping test execution due to a failure in the pod's events: %w", err)
				return
			}

		case <-timeoutChan:
			errs <- errors.New("time out")
			return
		}
	}
}

func waitForTaskRunPod(
	t *testing.T,
	c *kubernetes.Clientset,
	taskRunName,
	namespace string,
	errs chan error,
	taskRunDone chan bool,
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

	for {
		select {
		case err := <-errs:
			errs <- err
			return
		case <-taskRunDone:
			return
		case <-stop:
			podAdded <- taskRunPod
		}
	}
}
