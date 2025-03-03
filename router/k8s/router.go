package k8s

import (
	"github.com/gin-gonic/gin"
	"kubeants.com/api"
)

type K8SRouter struct{}

func (*K8SRouter) InitK8SRouter(r *gin.Engine) {
	group := r.Group("kapi")
	// 统一使用client-go的ForResource 方案转发k8s原生接口
	k8sResourceApiGroup := api.ApiGroupApp.K8sResourceApi
	// 动态获取有组名的资源
	group.Any("/apis/clusters/:cluster/:group/:version/namespaces/:namespace/:resource/*name", k8sResourceApiGroup.ProxyHandler)
	// 动态获取无组名的资源
	group.Any("/api/clusters/:cluster/:version/namespaces/:namespace/:resource/*name", k8sResourceApiGroup.ProxyHandler)

	// 下面是对每个k8s资源封装的接口，后期弃用
	// apiGroup := api.ApiGroupApp.K8SApiGroup
	// group.GET("pods", apiGroup.GetAllPods)
	// group.GET("pods/:namespace", apiGroup.GetPodsInNamespaceORDerail)
	// group.GET("pods/:namespace/:name", apiGroup.GetPodsInNamespaceORDerail)

	// group.POST("pods", apiGroup.CreateOrUpdatePod)
	// group.POST("pods/:namespace", apiGroup.CreateOrUpdatePod)
	// group.DELETE("pods/:namespace/:name", apiGroup.DeletePod)

	// group.GET("namespace", apiGroup.GetNamespaceList)
	// group.GET("namespace/:name", apiGroup.GetNamespaceList)

	// group.DELETE("namespace/:name", apiGroup.DeleteNamespace)

	// group.GET("node", apiGroup.GetNodeListOrDetail)
	// group.GET("node/:name", apiGroup.GetNodeListOrDetail)
}
