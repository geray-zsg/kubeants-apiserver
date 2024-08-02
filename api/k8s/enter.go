package k8s

import "kubeant.cn/service"

type K8SApi struct {
	PodApi
	NamespaceApi
	NodeApi
}

var podService = service.ServiceGroupApp.PodDerviceGroup.PodService
var namespaceService = service.ServiceGroupApp.NamespaceServiceGroup.NamespaceService
