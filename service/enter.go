package service

import (
	"kubeants.com/service/k8s"
	"kubeants.com/service/namespace"
	"kubeants.com/service/pod"
)

type ServiceGroup struct {
	PodDerviceGroup       pod.PodServiceGroup
	NamespaceServiceGroup namespace.NamespaceServiceGroup
	ResourceServiceGroup  k8s.GetResourcesGroup
}

var ServiceGroupApp = new(ServiceGroup)
