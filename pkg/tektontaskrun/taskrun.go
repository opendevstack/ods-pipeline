package tektontaskrun

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	k "github.com/opendevstack/ods-pipeline/internal/kubernetes"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelineclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/apis"
)

func TektonParamsFromStringParams(stringParams map[string]string) []tekton.Param {
	var params []tekton.Param
	for k, v := range stringParams {
		tp := tekton.Param{Name: k, Value: tekton.ParamValue{
			Type:      tekton.ParamTypeString,
			StringVal: v,
		}}
		params = append(params, tp)
	}
	return params
}

func runTask(tc *TaskRunConfig) (*tekton.TaskRun, bytes.Buffer, error) {
	clients := k.NewClients()
	tr, err := createTaskRunWithParams(clients.TektonClientSet, tc)
	if err != nil {
		return nil, bytes.Buffer{}, err
	}

	// TODO: if last output is short, it may be omitted from the logs.
	taskRun, logsBuffer, err := watchTaskRunUntilDone(clients, tc, tr)
	if err != nil {
		return nil, logsBuffer, err
	}

	log.Printf(
		"Task status: %q - %q\n",
		taskRun.Status.GetCondition(apis.ConditionSucceeded).GetReason(),
		taskRun.Status.GetCondition(apis.ConditionSucceeded).GetMessage(),
	)

	return taskRun, logsBuffer, nil
}

func createTaskRunWithParams(tknClient *pipelineclientset.Clientset, tc *TaskRunConfig) (*tekton.TaskRun, error) {

	taskWorkspaces := []tekton.WorkspaceBinding{}
	for wn, wd := range tc.Workspaces {
		if path.IsAbs(wd) && !strings.HasPrefix(wd, KinDMountHostPath) {
			return nil, fmt.Errorf("workspace dir %q is not located within %q", wd, KinDMountHostPath)
		}
		taskWorkspaces = append(taskWorkspaces, tekton.WorkspaceBinding{
			Name: wn,
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "task-pv-claim",
				ReadOnly:  false,
			},
			SubPath: strings.TrimPrefix(wd, KinDMountHostPath+"/"),
		})
	}

	tr, err := tknClient.TektonV1().TaskRuns(tc.Namespace).Create(context.TODO(),
		&tekton.TaskRun{
			ObjectMeta: metav1.ObjectMeta{
				Name: makeRandomTaskrunName(tc.Name),
			},
			Spec: tekton.TaskRunSpec{
				TaskRef:            &tekton.TaskRef{Kind: tekton.NamespacedTaskKind, Name: tc.Name},
				Params:             tc.Params,
				Workspaces:         taskWorkspaces,
				ServiceAccountName: tc.ServiceAccountName,
			},
		},
		metav1.CreateOptions{})

	return tr, err
}

func makeRandomTaskrunName(taskName string) string {
	return fmt.Sprintf("%s-taskrun-%s", taskName, makeRandomString(8))
}

func waitForTaskRunDone(
	ctx context.Context,
	c pipelineclientset.Interface,
	name, ns string,
	errs chan error,
	done chan *tekton.TaskRun) {

	deadline, _ := ctx.Deadline()
	timeout := time.Until(deadline)
	log.Printf("Waiting up to %v seconds for task %s in namespace %s to be done...\n", timeout.Round(time.Second).Seconds(), name, ns)

	w, err := c.TektonV1().TaskRuns(ns).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{
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
	podAdded chan *corev1.Pod) {
	log.Printf("Waiting for pod related to TaskRun %s to be added to the cluster\n", taskRunName)
	stop := make(chan struct{})

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(c, time.Second*30)
	podsInformer := kubeInformerFactory.Core().V1().Pods().Informer()

	var taskRunPod *corev1.Pod

	_, err := podsInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				// when a new task is created, watch its events
				pod := obj.(*corev1.Pod)
				if strings.HasPrefix(pod.Name, taskRunName) {
					taskRunPod = pod
					log.Printf("TaskRun %s added pod %s to the cluster", taskRunName, pod.Name)
					stop <- struct{}{}
				}

			},
		})
	if err != nil {
		log.Printf("Unable to install the event handler: %s", err)
	}

	defer close(stop)
	kubeInformerFactory.Start(stop)

	<-stop
	podAdded <- taskRunPod
}

func watchTaskRunUntilDone(c *k.Clients, tc *TaskRunConfig, tr *tekton.TaskRun) (*tekton.TaskRun, bytes.Buffer, error) {
	taskRunDone := make(chan *tekton.TaskRun)
	podAdded := make(chan *corev1.Pod)
	errs := make(chan error)
	collectedLogsChan := make(chan []byte)
	var collectedLogsBuffer bytes.Buffer

	ctx, cancel := context.WithTimeout(context.TODO(), tc.Timeout)
	defer cancel()
	go waitForTaskRunDone(
		ctx,
		c.TektonClientSet,
		tr.Name,
		tc.Namespace,
		errs,
		taskRunDone,
	)

	go waitForTaskRunPod(
		ctx,
		c.KubernetesClientSet,
		tr.Name,
		tc.Namespace,
		podAdded,
	)

	for {
		select {
		case err := <-errs:
			if err != nil {
				return nil, collectedLogsBuffer, err
			}

		case pod := <-podAdded:
			if pod != nil {
				go getEventsAndLogsOfPod(
					ctx,
					c.KubernetesClientSet,
					pod,
					collectedLogsChan,
					errs,
				)
			}

		case b := <-collectedLogsChan:
			collectedLogsBuffer.Write(b)

		case tr := <-taskRunDone:
			return tr, collectedLogsBuffer, nil
		case <-ctx.Done():
			return nil, collectedLogsBuffer, fmt.Errorf("timeout waiting for task run to finish. Consider increasing the timeout for your testcase at hand")
		}
	}
}
