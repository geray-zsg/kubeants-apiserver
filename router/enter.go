package router

import (
	"kubeants.io/router/example"
	"kubeants.io/router/k8s"
	"kubeants.io/router/kubeants"
)

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
	K8SRouterGroup     k8s.K8SRouter
	KaRouterGroup      kubeants.KaRouter
}

var RouterGroupApp = new(RouterGroup)
