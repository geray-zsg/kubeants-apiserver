package initiallize

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	kubeantconfig "kubeants.io/config"

	ctrl "sigs.k8s.io/controller-runtime"
)

func K8S() {
	var kubeconfig = "/root/.kube/config"

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// ✅ 设置 KubeRestConfig
	kubeantconfig.KubeRestConfig = config

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		ctrl.Log.V(2).Error(err, "创建客户端clientSet失败")
		panic(err.Error())
	}
	kubeantconfig.KubeClientSet = clientSet

	// create the dynamicClient
	dynameicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		ctrl.Log.V(2).Error(err, "创建动态客户端dynameicClient失败")
		panic(err.Error())
	}
	kubeantconfig.KubeDynamicClient = dynameicClient

	kubeantconfig.Kubeconfig = config.Host

}
