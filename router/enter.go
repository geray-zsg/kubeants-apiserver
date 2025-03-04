package router

import (
	"kubeants.io/router/example"
	"kubeants.io/router/k8s"
)

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
	K8SRouterGroup     k8s.K8SRouter
}

var RouterGroupApp = new(RouterGroup)
