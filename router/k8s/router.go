package k8s

import (
	"github.com/gin-gonic/gin"
	"kubeants.io/api"
	"kubeants.io/middleware"
)

type K8SRouter struct{}

func (*K8SRouter) InitK8SRouter(r *gin.Engine) {
	group := r.Group("/gapi/cluster/:cluster/workspace/:workspace")
	// 统一使用client-go的ForResource 方案转发k8s原生接口
	k8sResourceApiGroup := api.ApiGroupApp.K8sResourceApi
	// 认证接口
	group.Use(middleware.AuthMiddleware())

	// 动态获取无组名的资源
	group.Any("/api/:version/:resource/*name", k8sResourceApiGroup.ProxyHandler)
	// 无组名namespace级别，例如Pod
	group.Any("/api/:version/namespaces/:namespace/:resource/*name", k8sResourceApiGroup.ProxyHandler)
	// 动态获取有组名的资源
	group.Any("/apis/:group/:version/namespaces/:namespace/:resource/*name", k8sResourceApiGroup.ProxyHandler)
	// 集群界别资源有组名，无需单独提供接口，后端config.KubeDynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})  可根据提供的namespace为""自动过滤掉，例如：clusterrole
	group.Any("/apis/:group/:version/:resource/*name", k8sResourceApiGroup.ProxyHandler)

}
