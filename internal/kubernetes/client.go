package kubernetes

import (
	"flag"
	"path/filepath"

	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Clients struct {
	KubernetesClientSet *kubernetes.Clientset
	TektonClientSet     *tekton.Clientset
}

func NewClients() *Clients {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
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
