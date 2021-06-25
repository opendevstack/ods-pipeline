// from https://github.com/tektoncd/cli/blob/c996b3004650658c73ae8b1c0ce98f5165107af5/test/framework/logs.go
package tasktesting

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"knative.dev/pkg/test/logging"
)

// CollectPodLogs will get the logs for all containers in a Pod
func CollectPodLogs(c *kubernetes.Clientset, podName, namespace string, logf logging.FormatLogger, podEventsDone chan<- bool) {
	logs, err := getContainerLogsFromPod(c, podName, namespace, podEventsDone)
	if err != nil {
		logf("Could not get logs for pod %s: %s", podName, err)
	}
	logf("pod logs %s", logs)
}

func getContainerLogsFromPod(c kubernetes.Interface, pod, namespace string, podEventsDone chan<- bool) (string, error) {
	p, err := c.CoreV1().Pods(namespace).Get(context.Background(), pod, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("could not get pod %s in namespace %s: %w", pod, namespace, err)
	}

	sb := strings.Builder{}
	for _, container := range p.Spec.Containers {
		log.Printf("Waiting for container %s from pod %s to be Ready...\n", container.Name, pod)

		deadline := time.Now().Add(30 * time.Second)
		for {
			containerIsReady := false
			for _, cs := range p.Status.ContainerStatuses {
				if cs.Name == container.Name && cs.Ready {
					containerIsReady = true
				}
			}
			if containerIsReady || time.Now().After(deadline) {
				break
			}
			time.Sleep(3 * time.Second)
		}

		sb.WriteString(fmt.Sprintf("\n>>> Container %s:\n", container.Name))
		req := c.CoreV1().Pods(namespace).GetLogs(pod, &corev1.PodLogOptions{
			Follow:    true,
			Container: container.Name,
		})
		rc, err := req.Stream(context.Background())
		if err != nil {
			return "", fmt.Errorf("could not create log stream for pod %s in namespace %s: %w", pod, namespace, err)
		}
		bs, err := ioutil.ReadAll(rc)
		if err != nil {
			return "", fmt.Errorf("could not read log stream for pod %s in namespace %s: %w", pod, namespace, err)
		}
		sb.Write(bs)
	}

	podEventsDone <- true // notify the channel to stop watching pod events

	return sb.String(), nil // display containers' logs
}
