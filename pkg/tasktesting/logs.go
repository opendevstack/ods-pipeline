package tasktesting

import (
	"context"
	"fmt"
	"io/ioutil"
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
	ctx context.Context,
	c kubernetes.Interface,
	pod *corev1.Pod,
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

	for _, container := range pod.Spec.Containers {
		err := streamContainerLogs(ctx, c, podNamespace, podName, container.Name)
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
// TODO: Make this watch the pod.
// When the container is not waiting anymore, start logs. when thas has been done
// once, then in the next loop if the state is terminated, we stop the logs.
// when logs have been stopped, do not block anymore and go to next container.
func waitForContainerReady(
	ctx context.Context,
	c kubernetes.Interface,
	podNamespace, podName, containerName string) error {
	ticker := time.NewTicker(2 * time.Second)
	deadline := time.Now().Add(5 * time.Minute)
	for {
		<-ticker.C
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for container %s to become ready", containerName)
		}
		p, err := c.CoreV1().Pods(podNamespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("could not get pod %s in namespace %s: %w", podName, podNamespace, err)
		}
		for _, cs := range p.Status.ContainerStatuses {
			if cs.Name == containerName {
				if cs.State.Running != nil || cs.State.Terminated != nil {
					log.Printf("Container %s is ready", containerName)
					return nil
				}
			}
		}
	}
}

// streamContainerLogs waits for container to be ready, then streams the logs.
func streamContainerLogs(
	ctx context.Context,
	c kubernetes.Interface,
	podNamespace, podName, containerName string) error {
	log.Printf("Waiting for container %s from pod %s to be ready...\n", containerName, podName)

	err := waitForContainerReady(ctx, c, podNamespace, podName, containerName)
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

	logs, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("could not read log stream for pod %s in namespace %s: %w", podName, podNamespace, err)
	}
	fmt.Println(string(logs))
	// TODO: stream the logs as they come. For this we need to figure out when to
	// start and when to stop streaming.
	// for {
	// 	buf := make([]byte, 100)

	// 	numBytes, err := rc.Read(buf)
	// 	if numBytes == 0 {
	// 		continue
	// 	}
	// 	if err == io.EOF {
	// 		fmt.Printf("logs for %s ended\n", containerName)
	// 		break
	// 	}
	// 	if err != nil {
	// 		return fmt.Errorf("error in copy information from podLogs to buf: %w", err)
	// 	}

	// 	fmt.Print(string(buf[:numBytes]))
	// }
	return nil
}
