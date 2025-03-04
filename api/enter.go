package api

import (
	"kubeants.io/api/example"
	"kubeants.io/api/k8s"
)

type ApiGroup struct {
	ExampleApiGroup example.ExampleTestApi
	K8SApiGroup     k8s.K8SApi
	K8sResourceApi  k8s.ResourceApi // 所有k8s资源的API
}

var ApiGroupApp = new(ApiGroup)
