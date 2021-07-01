package kubernetes

import (
	"context"
	"fmt"
	"log"

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
