package tasktesting

import (
	"fmt"
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WaitForTaskRunPod(t *testing.T, c *kubernetes.Clientset, taskRunName, namespace string) *v1.Pod {
	fmt.Println("waiting for pod related to taskrun", taskRunName)
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
					fmt.Println("found pod", pod.Name)
					stop <- struct{}{}
					//WatchPodEvents(t, c, pod.Name, namespace, podEventsDone)
				}

			},
		})

	defer close(stop)
	kubeInformerFactory.Start(stop)

	for {
		<-stop
		return taskRunPod
	}
}
