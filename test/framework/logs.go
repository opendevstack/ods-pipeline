// from https://github.com/tektoncd/cli/blob/c996b3004650658c73ae8b1c0ce98f5165107af5/test/framework/logs.go
package framework

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"knative.dev/pkg/test/logging"
)

// CollectPodLogs will get the logs for all containers in a Pod
func CollectPodLogs(c *kubernetes.Clientset, podName, namespace string, logf logging.FormatLogger) {
	logs, err := getContainerLogsFromPod(c, podName, namespace)
	if err != nil {
		logf("Could not get logs for pod %s: %s", podName, err)
	}
	logf("pod logs %s", logs)
}

func getContainerLogsFromPod(c kubernetes.Interface, pod, namespace string) (string, error) {
	p, err := c.CoreV1().Pods(namespace).Get(context.Background(), pod, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	for _, container := range p.Spec.Containers {
		sb.WriteString(fmt.Sprintf("\n>>> Container %s:\n", container.Name))
		req := c.CoreV1().Pods(namespace).GetLogs(pod, &corev1.PodLogOptions{Follow: true, Container: container.Name})
		rc, err := req.Stream(context.Background())
		if err != nil {
			return "", err
		}
		bs, err := ioutil.ReadAll(rc)
		if err != nil {
			return "", err
		}
		sb.Write(bs)
	}
	return sb.String(), nil
}
