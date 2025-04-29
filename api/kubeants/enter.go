package kubeants

import "kubeants.io/service"

type KubeantsApi struct {
	UserApi
}

var resourceService = service.ServiceGroupApp.ResourceServiceGroup.ResourceService
