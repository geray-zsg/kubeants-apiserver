package k8s

import (
	"github.com/gin-gonic/gin"
	"kubeant.cn/api"
)

type K8SRouter struct{}

func (*K8SRouter) InitK8SRouter(r *gin.Engine) {
	group := r.Group("kapi")
	apiGroup := api.ApiGroupApp.K8SApiGroup
	group.GET("pods", apiGroup.GetAllPods)
	group.GET("pods/:namespace", apiGroup.GetPodsInNamespaceORDerail)
	group.GET("pods/:namespace/:name", apiGroup.GetPodsInNamespaceORDerail)

	group.POST("pods", apiGroup.CreateOrUpdatePod)
	group.POST("pods/:namespace", apiGroup.CreateOrUpdatePod)
	group.DELETE("pods/:namespace/:name", apiGroup.DeletePod)

	group.GET("namespace", apiGroup.GetNamespaceList)
	group.GET("namespace/:name", apiGroup.GetNamespaceList)

	group.DELETE("namespace/:name", apiGroup.DeleteNamespace)

	group.GET("node", apiGroup.GetNodeListOrDetail)
	group.GET("node/:name", apiGroup.GetNodeListOrDetail)
}
