package service

import (
	"kubeant.cn/service/namespace"
	"kubeant.cn/service/pod"
)

type ServiceGroup struct {
	PodDerviceGroup       pod.PodServiceGroup
	NamespaceServiceGroup namespace.NamespaceServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
