package service

import (
	"kubeants.io/service/k8s"
	"kubeants.io/service/namespace"
	"kubeants.io/service/pod"
)

type ServiceGroup struct {
	PodDerviceGroup       pod.PodServiceGroup
	NamespaceServiceGroup namespace.NamespaceServiceGroup
	ResourceServiceGroup  k8s.GetResourcesGroup
	LogService            k8s.LogService
}

var ServiceGroupApp = new(ServiceGroup)
