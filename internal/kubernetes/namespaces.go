package kubernetes

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateRandomNamespace(clientset *kubernetes.Clientset) string {

	id := uuid.NewV4()

	ns, err := clientset.CoreV1().Namespaces().Create(context.TODO(),
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: id.String(),
			},
		},
		metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created Namespace: %s\n", ns.Name)

	return ns.Name
}
