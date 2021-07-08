package tasktesting

import (
	"context"
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func watchPodEvents(
	ctx context.Context,
	c kubernetes.Interface,
	podName, namespace string,
	stop chan bool,
	errs chan error) {

	log.Printf("Watching events for pod %s in namespace %s", podName, namespace)

	ew, err := c.CoreV1().Events(namespace).Watch(context.Background(),
		metav1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s", podName, namespace),
		})
	if err != nil {
		errs <- fmt.Errorf("failed to watch events from pod %s in namespace %s", podName, namespace)
		return
	}

	log.Println("---------------------- Events -------------------------")

	// Wait for any event failure or a all its containers to be running
	for {
		select {
		case wev := <-ew.ResultChan():
			if wev.Object != nil {
				ev := wev.Object.(*v1.Event)
				log.Printf("Type: %s, Message: %s", ev.Type, ev.Message)
				if ev.Type == "Warning" && strings.Contains(ev.Message, "Error") {
					errs <- fmt.Errorf("error detected in events: %s", ev.Message)
					return
				}
			}
		case <-stop:
			fmt.Println("quit watching events as no more are expected")
			return
		}
	}
}
