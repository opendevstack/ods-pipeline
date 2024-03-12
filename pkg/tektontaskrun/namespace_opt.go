package tektontaskrun

import (
	"context"
	"log"

	k "github.com/opendevstack/ods-pipeline/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NamespaceOpt allows to further configure the K8s namespace after its creation.
type NamespaceOpt func(cc *ClusterConfig, nc *NamespaceConfig) error

// NamespaceConfig represents key configuration of the K8s namespace.
type NamespaceConfig struct {
	Name string
}

// SetupTempNamespace sets up a new namespace using a pseduo-random name,
// applies any given NamespaceOpt and returns a function to clean up the
// namespace at a later time.
func SetupTempNamespace(cc *ClusterConfig, opts ...NamespaceOpt) (nc *NamespaceConfig, cleanup func(), err error) {
	nc = &NamespaceConfig{
		Name: makeRandomString(8),
	}
	cleanup, err = initNamespaceAndPVC(cc, nc)
	if err != nil {
		return
	}
	cleanupOnInterrupt(cleanup)
	for _, o := range opts {
		err = o(cc, nc)
		if err != nil {
			return
		}
	}
	return
}

// InstallTaskFromPath renders the task template at path using the given data,
// then installs the resulting task into the namespace identified by
// NamespaceConfig.
func InstallTaskFromPath(path string, data map[string]string) NamespaceOpt {
	return func(cc *ClusterConfig, nc *NamespaceConfig) error {
		d := cc.DefaultManifestTemplateData()
		for k, v := range data {
			d[k] = v
		}
		_, err := installTask(path, nc.Name, d)
		return err
	}
}

func initNamespaceAndPVC(cc *ClusterConfig, nc *NamespaceConfig) (cleanup func(), err error) {
	clients := k.NewClients()

	_, nsCleanup, err := createTempNamespace(clients.KubernetesClientSet, nc.Name)
	if err != nil {
		return nil, err
	}

	// for simplicity and traceability, use namespace name for PVC as well
	_, pvcCleanup, err := createTempPVC(clients.KubernetesClientSet, cc, nc.Name)
	if err != nil {
		return nil, err
	}

	return func() {
		nsCleanup()
		pvcCleanup()
	}, nil
}

func createTempNamespace(clientset kubernetes.Interface, name string) (namespace *corev1.Namespace, cleanup func(), err error) {
	namespace, err = clientset.CoreV1().Namespaces().Create(
		context.TODO(),
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		},
		metav1.CreateOptions{},
	)
	return namespace, func() {
		log.Printf("Removing temporary namespace %q ...", name)
		err := removeNamespace(clientset, name)
		if err != nil {
			log.Println(err)
		}
	}, err
}

func createTempPVC(clientset kubernetes.Interface, cc *ClusterConfig, name string) (pvc *corev1.PersistentVolumeClaim, cleanup func(), err error) {
	_, err = k.CreatePersistentVolume(
		clientset,
		name,
		cc.StorageCapacity,
		cc.StorageSourceDir,
		cc.StorageClassName,
	)
	if err != nil {
		return
	}

	pvc, err = k.CreatePersistentVolumeClaim(
		clientset,
		cc.StorageCapacity,
		cc.StorageClassName,
		name,
	)
	if err != nil {
		return
	}
	return pvc, func() {
		err := removePVC(clientset, name)
		if err != nil {
			log.Println(err)
		}
	}, err
}

func removeNamespace(clientset kubernetes.Interface, name string) error {
	return clientset.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
}

func removePVC(clientset kubernetes.Interface, name string) error {
	return clientset.CoreV1().PersistentVolumes().Delete(context.Background(), name, metav1.DeleteOptions{})
}
