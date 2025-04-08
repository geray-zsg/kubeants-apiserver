package k8s

import "kubeants.io/service"

type K8SApi struct {
	ResourceApi
}

var resourceService = service.ServiceGroupApp.ResourceServiceGroup.ResourceService // 所有k8s资源原生接口
