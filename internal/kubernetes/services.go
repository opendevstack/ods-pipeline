package kubernetes

import (
	"context"
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateNodePortService(clientset *kubernetes.Clientset, name string, selectors map[string]string, port, targetPort int32, namespace string) (*v1.Service, error) {

	log.Printf("Create node port service %s", name)
	svc, err := clientset.CoreV1().Services(namespace).Create(context.TODO(),
		&v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:   name,
				Labels: map[string]string{"app.kubernetes.io/managed-by": "ods-pipeline"},
			},
			Spec: v1.ServiceSpec{
				ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
				Ports: []v1.ServicePort{
					{
						Name:       fmt.Sprintf("%d-%d", port, targetPort),
						NodePort:   port,
						Port:       port,
						Protocol:   v1.ProtocolTCP,
						TargetPort: intstr.FromInt(int(targetPort)),
					},
				},
				Selector: map[string]string{
					"eventlistener": "ods-pipeline",
				},
				SessionAffinity: v1.ServiceAffinityNone,
				Type:            v1.ServiceTypeNodePort,
			},
		}, metav1.CreateOptions{})

	return svc, err
}

// ServiceHasReadyPods returns false if no pod is assigned to given service
// or if one or more pods are not "Running"
// or one or more of any pods containers are not "ready".
func ServiceHasReadyPods(clientset *kubernetes.Clientset, svc *v1.Service) (bool, string, error) {
	podList, err := servicePods(clientset, svc)
	if err != nil {
		return false, "error", err
	}
	for _, pod := range podList.Items {
		phase := pod.Status.Phase
		if phase != "Running" {
			return false, fmt.Sprintf("pod %s is in phase %+v", pod.Name, phase), nil
		}
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if !containerStatus.Ready {
				return false, fmt.Sprintf("container %s in pod %s is not ready", containerStatus.Name, pod.Name), nil
			}
		}
	}
	return true, "ok", nil
}

func servicePods(clientset *kubernetes.Clientset, svc *v1.Service) (*v1.PodList, error) {
	podClient := clientset.CoreV1().Pods(svc.Namespace)
	selector := []string{}
	for key, value := range svc.Spec.Selector {
		selector = append(selector, fmt.Sprintf("%s=%s", key, value))
	}
	pods, err := podClient.List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: strings.Join(selector, ","),
		},
	)
	if err != nil {
		return nil, err
	}
	return pods.DeepCopy(), nil
}
