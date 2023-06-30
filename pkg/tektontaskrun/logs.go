package tektontaskrun

import (
	"bufio"
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// getEventsAndLogsOfPod streams events of the pod until all containers are ready,
// and streams logs for each container once ready. It stops if there are any
// sends on the errs channels or if the passed context is cancelled.
func getEventsAndLogsOfPod(
	ctx context.Context,
	c kubernetes.Interface,
	pod *corev1.Pod,
	collectedLogsChan chan []byte,
	errs chan error) {
	quitEvents := make(chan bool)
	podName := pod.Name
	podNamespace := pod.Namespace

	go watchPodEvents(
		ctx,
		c,
		podName,
		podNamespace,
		quitEvents,
		errs,
	)

	watchingEvents := true
	for _, container := range pod.Spec.Containers {
		err := streamContainerLogs(ctx, c, podNamespace, podName, container.Name, collectedLogsChan)
		if err != nil {
			fmt.Printf("failure while getting container logs: %s", err)
			errs <- err
			return
		}
		if watchingEvents {
			quitEvents <- true
			watchingEvents = false
		}
	}
}

func streamContainerLogs(
	ctx context.Context,
	c kubernetes.Interface,
	podNamespace, podName, containerName string, collectedLogsChan chan []byte) error {
	log.Printf("Waiting for container %s from pod %s to be ready...\n", containerName, podName)

	w, err := c.CoreV1().Pods(podNamespace).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{
		Name:      podName,
		Namespace: podNamespace,
	}))
	if err != nil {
		return fmt.Errorf("error watching pods: %s", err)
	}

	for {
		ev := <-w.ResultChan()
		if cs, ok := containerFromEvent(ev, podName, containerName); ok {
			if cs.State.Running != nil {
				log.Printf("---------------------- Logs from %s -------------------------\n", containerName)
				// Set up log stream using a new ctx so that it's not cancelled
				// when the task is done before all logs have been read.
				ls, err := c.CoreV1().Pods(podNamespace).GetLogs(podName, &corev1.PodLogOptions{
					Follow:    true,
					Container: containerName,
				}).Stream(context.Background())
				if err != nil {
					return fmt.Errorf("could not create log stream for pod %s in namespace %s: %w", podName, podNamespace, err)
				}
				defer ls.Close()
				reader := bufio.NewScanner(ls)
				for reader.Scan() {
					select {
					case <-ctx.Done():
						collectedLogsChan <- reader.Bytes()
						fmt.Println(reader.Text())
						return nil
					default:
						collectedLogsChan <- reader.Bytes()
						fmt.Println(reader.Text())
					}
				}
				return reader.Err()
			}
		}
	}
}

func containerFromEvent(ev watch.Event, podName, containerName string) (corev1.ContainerStatus, bool) {
	if ev.Object != nil {
		p, ok := ev.Object.(*corev1.Pod)
		if ok && p.Name == podName {
			for _, cs := range p.Status.ContainerStatuses {
				if cs.Name == containerName {
					return cs, true
				}
			}
		}
	}
	return corev1.ContainerStatus{}, false
}
