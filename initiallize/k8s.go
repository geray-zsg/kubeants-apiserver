package initiallize

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	kubeantconfig "kubeants.io/config"
)

func K8S() {
	var kubeconfig = "/root/.kube/config"

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("kubeconfig apiserverHost:", config.Host)
	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create clientset: %v", err)
		panic(err.Error())
	}
	kubeantconfig.KubeClientSet = clientSet

	// create the dynamicClient
	dynameicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create dynamicClient: %v", err)
		panic(err.Error())
	}
	kubeantconfig.KubeDynamicClient = dynameicClient

	kubeantconfig.Kubeconfig = config.Host

}
