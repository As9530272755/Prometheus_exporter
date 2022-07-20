package linkKube

import (
	"kube_expoter/msgErr"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Link_kube(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	msgErr.ErrInfo(err)

	clientSet, err := kubernetes.NewForConfig(config)
	msgErr.ErrInfo(err)

	return clientSet
}
