package tasktesting

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WatchTaskRunEvents(c *kubernetes.Clientset, taskRunName, namespace string) {

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(c, time.Second*30)
	podsInformer := kubeInformerFactory.Core().V1().Pods().Informer()

	podsInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				// when a new task is created, watch its events
				pod := obj.(*v1.Pod)
				if strings.HasPrefix(pod.Name, taskRunName) {
					WatchPodEvents(c, pod.Name, namespace)
				}
			},
		})

	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)
	for {
		time.Sleep(time.Second)
	}
}

func WatchPodEvents(c *kubernetes.Clientset, podName, namespace string) {

	log.Printf("Watching events for pod %s in namespace %s", podName, namespace)

	time.Sleep(3 * time.Second)

	events, err := c.CoreV1().Events(namespace).List(context.TODO(),
		metav1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s", podName, namespace),
		})
	if err != nil {
		log.Fatalf("Failed to display events from pod %s in namespace %s\n", podName, namespace)
	}

	log.Println("------------- Events ------------------")
	for _, e := range events.Items {
		log.Printf("Type: %s, Message: %s", e.Type, e.Message)
	}
	log.Println("-------------------------------------------------")
}
