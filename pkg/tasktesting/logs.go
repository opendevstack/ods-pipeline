package tasktesting

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// getEventsAndLogsOfPod streams events of the pod until all containers are ready,
// and streams logs for each container once ready. It stops if there are any
// sends on the errs or taskRunDone channels.
func getEventsAndLogsOfPod(
	c kubernetes.Interface,
	pod *corev1.Pod,
	errs chan error,
	taskRunDone chan bool) {
	quitEvents := make(chan bool)
	podName := pod.Name
	podNamespace := pod.Namespace

	go watchPodEvents(
		c,
		podName,
		podNamespace,
		quitEvents,
		errs,
		taskRunDone,
	)

	for _, container := range pod.Spec.Containers {
		err := streamContainerLogs(c, podNamespace, podName, container.Name, errs, taskRunDone)
		if err != nil {
			fmt.Println("failure while getting container logs")
			errs <- err
			return
		}
	}
	fmt.Println("done with the logs, quitting events")
	quitEvents <- true
}

// waitForContainerReady waits until the container is "Ready" for up to 5 minutes.
func waitForContainerReady(
	c kubernetes.Interface,
	podNamespace, podName, containerName string,
	errs chan error,
	taskRunDone chan bool) error {
	ticker := time.NewTicker(2 * time.Second)
	deadline := time.Now().Add(5 * time.Minute)
	for {
		select {
		case <-taskRunDone:
			return nil
		case err := <-errs:
			return err
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timed out waiting for container %s to become ready", containerName)
			}
			p, err := c.CoreV1().Pods(podNamespace).Get(context.Background(), podName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("could not get pod %s in namespace %s: %w", podName, podNamespace, err)
			}
			for _, cs := range p.Status.ContainerStatuses {
				if cs.Name == containerName && cs.Ready {
					log.Printf("Container %s is ready", containerName)
					return nil
				}
			}
		}
	}
}

// streamContainerLogs waits for container to be ready, then streams the logs.
func streamContainerLogs(
	c kubernetes.Interface,
	podNamespace, podName, containerName string,
	errs chan error,
	taskRunDone chan bool) error {
	log.Printf("Waiting for container %s from pod %s to be ready...\n", containerName, podName)

	err := waitForContainerReady(c, podNamespace, podName, containerName, errs, taskRunDone)
	if err != nil {
		return err
	}

	log.Printf(">>> Container %s:\n", containerName)
	req := c.CoreV1().Pods(podNamespace).GetLogs(podName, &corev1.PodLogOptions{
		Follow:    true,
		Container: containerName,
	})
	rc, err := req.Stream(context.Background())
	if err != nil {
		return fmt.Errorf("could not create log stream for pod %s in namespace %s: %w", podName, podNamespace, err)
	}
	defer rc.Close()
	for {
		buf := make([]byte, 100)

		numBytes, err := rc.Read(buf)
		if numBytes == 0 {
			continue
		}
		if err == io.EOF {
			fmt.Printf("logs for %s ended\n", containerName)
			break
		}
		if err != nil {
			return fmt.Errorf("error in copy information from podLogs to buf: %w", err)
		}

		fmt.Print(string(buf[:numBytes]))
	}
	return nil
}
