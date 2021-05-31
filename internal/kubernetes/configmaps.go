package kubernetes

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetConfigMap(clientset *kubernetes.Clientset, namespace string, cmName string) (*v1.ConfigMap, error) {

	log.Printf("Get configmap %s", cmName)

	cm, err := clientset.CoreV1().
		ConfigMaps(namespace).
		Get(context.TODO(), cmName, metav1.GetOptions{})

	return cm, err
}

func GetConfigMapKey(clientset *kubernetes.Clientset, namespace, cmName, key string) (string, error) {

	log.Printf("Get configmap %s", cmName)

	cm, err := clientset.CoreV1().
		ConfigMaps(namespace).
		Get(context.TODO(), cmName, metav1.GetOptions{})

	if err != nil {
		return "", err
	}

	v, ok := cm.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	return v, err
}
