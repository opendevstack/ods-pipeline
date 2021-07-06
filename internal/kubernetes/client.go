package kubernetes

import (
	"path/filepath"

	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Clients struct {
	KubernetesClientSet *kubernetes.Clientset
	TektonClientSet     *tekton.Clientset
}

func NewClients() *Clients {
	// TODO: make configurable from outside
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the Kubernetes clientset
	kubernetesClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create the Tekton clientset
	tektonClientSet, err := tekton.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Clients{
		KubernetesClientSet: kubernetesClientset,
		TektonClientSet:     tektonClientSet,
	}
}

func NewInClusterClientset() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	return kubernetes.NewForConfig(config)
}
