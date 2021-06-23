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

func WatchTaskRunEvents(c *kubernetes.Clientset, taskRunName, namespace string, timeout time.Duration) {

	stop := make(chan struct{})

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(c, time.Second*30)
	podsInformer := kubeInformerFactory.Core().V1().Pods().Informer()

	podsInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				// when a new task is created, watch its events
				pod := obj.(*v1.Pod)
				if strings.HasPrefix(pod.Name, taskRunName) {
					WatchPodEvents(c, pod.Name, namespace, timeout)
					stop <- struct{}{}
				}

			},
		})

	defer close(stop)
	kubeInformerFactory.Start(stop)

	for {
		select {
		case <-stop:
			return
		}
	}
}

func WatchPodEvents(c *kubernetes.Clientset, podName, namespace string, timeout time.Duration) {

	log.Printf("Watching events for pod %s in namespace %s", podName, namespace)

	time.Sleep(3 * time.Second) //TODO: How to wait until Pod is actually created?

	ew, err := c.CoreV1().Events(namespace).Watch(context.Background(),
		metav1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s", podName, namespace),
		})
	if err != nil {
		log.Fatalf("Failed to watch events from pod %s in namespace %s\n", podName, namespace)
	}

	// Setup a timeout channel
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		timeoutChan <- struct{}{}
	}()

	log.Println("---------------------- Events -------------------------")

	// Wait for any failure or a timeout
	for {
		select {
		case wev := <-ew.ResultChan():
			if wev.Object != nil {
				ev := wev.Object.(*v1.Event)
				log.Printf("Type: %s, Message: %s", ev.Type, ev.Message)
				if ev.Type == "Warning" && strings.Contains(ev.Message, "Error") {
					log.Fatalf("The following error has been detected in the events output: %s\n", ev.Message)
					//TODO: When it fails we have to clean up the namespace, pvc, etc...
					break
				}
			}
		case <-timeoutChan:
			log.Println("-----------------------------------------------")
			log.Printf("No failures detected in the events output of pod %s after %v seconds\n", podName, timeout.Seconds())
			return
		}

	}

}
