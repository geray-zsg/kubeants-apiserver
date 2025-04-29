package api

import (
	"kubeants.io/api/example"
	"kubeants.io/api/k8s"
	"kubeants.io/api/kubeants"
)

type ApiGroup struct {
	ExampleApiGroup example.ExampleTestApi
	K8SApiGroup     k8s.K8SApi
	K8sResourceApi  k8s.ResourceApi      // 所有k8s资源的API
	KaApi           kubeants.KubeantsApi // kubeants自定义资源接口
}

var ApiGroupApp = new(ApiGroup)
