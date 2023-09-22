package kubernetes

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// TODO: return error
func CreateNamespace(clientset kubernetes.Interface, namespace string) {
	log.Printf("Create namespace %s to deploy to", namespace)
	if _, err := clientset.CoreV1().Namespaces().Create(context.TODO(),
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		},
		metav1.CreateOptions{}); err != nil {
		log.Printf("Failed to create namespace %s for tests: %s", namespace, err)
	}
}
