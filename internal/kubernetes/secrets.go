package kubernetes

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetSecret(clientset *kubernetes.Clientset, namespace string, secretName string) (*v1.Secret, error) {

	log.Printf("Get secret %s", secretName)

	secret, err := clientset.CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secretName, metav1.GetOptions{})

	return secret, err
}

func GetSecretKey(clientset *kubernetes.Clientset, namespace, secretName, key string) (string, error) {

	log.Printf("Get secret %s", secretName)

	secret, err := clientset.CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secretName, metav1.GetOptions{})

	if err != nil {
		return "", err
	}

	v, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	return string(v), err
}
