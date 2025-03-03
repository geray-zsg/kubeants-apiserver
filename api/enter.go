package api

import (
	"kubeants.com/api/example"
	"kubeants.com/api/k8s"
)

type ApiGroup struct {
	ExampleApiGroup example.ExampleTestApi
	K8SApiGroup     k8s.K8SApi
}

var ApiGroupApp = new(ApiGroup)
