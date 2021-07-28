package tasktesting

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

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
	errs chan error,
	podLogs chan []byte) {
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
		err := streamContainerLogs(ctx, c, podNamespace, podName, container.Name, podLogs)
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
	podNamespace, podName, containerName string,
	podLogs chan []byte) error {
	log.Printf("Waiting for container %s from pod %s to be ready...\n", containerName, podName)

	w, err := c.CoreV1().Pods(podNamespace).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{
		Name:      podName,
		Namespace: podNamespace,
	}))
	if err != nil {
		return fmt.Errorf("error watching pods: %s", err)
	}

	containerState := "waiting"
	var logStream io.ReadCloser
	for {
		select {
		case ev := <-w.ResultChan():
			if cs, ok := containerFromEvent(ev, podName, containerName); ok {
				if cs.State.Running != nil {
					if containerState == "waiting" {
						log.Printf("---------------------- Logs from %s -------------------------\n", containerName)
						req := c.CoreV1().Pods(podNamespace).GetLogs(podName, &corev1.PodLogOptions{
							Follow:    true,
							Container: containerName,
						})
						ls, err := req.Stream(context.Background())
						if err != nil {
							return fmt.Errorf("could not create log stream for pod %s in namespace %s: %w", podName, podNamespace, err)
						}
						logStream = ls
						defer logStream.Close()
					}
					containerState = "running"
				}
				if containerState != "waiting" && cs.State.Terminated != nil {
					// read reminder of the log stream
					logs, err := ioutil.ReadAll(logStream)
					if err != nil {
						return fmt.Errorf("could not read log stream for pod %s in namespace %s: %w", podName, podNamespace, err)
					}
					fmt.Println(string(logs))
					return nil
				}
			}

		default:
			// if log stream has started, read some bytes
			if logStream != nil {
				buf := make([]byte, 100)

				numBytes, err := logStream.Read(buf)
				if numBytes == 0 {
					continue
				}
				if err == io.EOF {
					log.Printf("logs for %s ended\n", containerName)
					return nil
				}
				if err != nil {
					return fmt.Errorf("error in reading log stream: %w", err)
				}

				b := buf[:numBytes]
				podLogs <- b
				fmt.Print(string(b))
			} else {
				time.Sleep(time.Second)
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
