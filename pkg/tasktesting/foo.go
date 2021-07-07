package tasktesting

import (
	"log"
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

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
