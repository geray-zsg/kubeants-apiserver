package api

import (
	"kubeant.cn/api/example"
	"kubeant.cn/api/k8s"
)

type ApiGroup struct {
	ExampleApiGroup example.ExampleTestApi
	K8SApiGroup     k8s.K8SApi
}

var ApiGroupApp = new(ApiGroup)
