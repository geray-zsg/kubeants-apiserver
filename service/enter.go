package service

import (
	"kubeants.com/service/namespace"
	"kubeants.com/service/pod"
)

type ServiceGroup struct {
	PodDerviceGroup       pod.PodServiceGroup
	NamespaceServiceGroup namespace.NamespaceServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
