// router/k8s/router.go
package k8s

import (
	"github.com/gin-gonic/gin"
	"kubeants.io/api"
	"kubeants.io/middleware"
)

type K8SRouter struct{}

func (*K8SRouter) InitK8SRouter(r *gin.Engine) {

	// 统一使用client-go的ForResource 方案转发k8s原生接口
	k8sResourceApiGroup := api.ApiGroupApp.K8sResourceApi
	k8sLogApiGroup := api.ApiGroupApp.LogApi
	// ========== 集群级资源：cluster-scoped ==========
	cluster := r.Group("/gapi/cluster/:cluster")
	cluster.Use(middleware.AuthMiddleware())

	cluster.Any("/api/:version/:resource", k8sResourceApiGroup.ProxyHandler)
	cluster.Any("/api/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)

	cluster.Any("/apis/:group/:version/:resource", k8sResourceApiGroup.ProxyHandler)
	cluster.Any("/apis/:group/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)

	// ========== 工作空间内的命名空间资源 ==========
	workspace := r.Group("/gapi/cluster/:cluster/workspace/:workspace")
	workspace.Use(middleware.AuthMiddleware())

	workspace.Any("/api/:version/namespaces/:namespace/:resource", k8sResourceApiGroup.ProxyHandler)
	workspace.Any("/api/:version/namespaces/:namespace/:resource/:name", k8sResourceApiGroup.ProxyHandler)

	workspace.Any("/apis/:group/:version/namespaces/:namespace/:resource", k8sResourceApiGroup.ProxyHandler)
	workspace.Any("/apis/:group/:version/namespaces/:namespace/:resource/:name", k8sResourceApiGroup.ProxyHandler)

	// 补充，如果数据 属于 cluster-scope 类型，则走cluster路由，同时不要影响之前的前端接口逻辑（保留部分）
	workspace.Any("/api/:version/:resource", k8sResourceApiGroup.ProxyHandler)
	workspace.Any("/api/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)
	workspace.Any("/apis/:group/:version/:resource", k8sResourceApiGroup.ProxyHandler)
	workspace.Any("/apis/:group/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)

	// Pod日志接口
	workspace.GET("/api/v1/namespaces/:namespace/pods/:pod/log", k8sLogApiGroup.GetPodLogs)

	// 认证接口
	// r.Use(middleware.AuthMiddleware())
	// // r.Any("/gapi/cluster/:cluster/workspace/:workspace/api/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)
	// // r.Any("/gapi/cluster/:cluster/workspace/:workspace/apis/:group/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)
	// group := r.Group("/gapi/cluster/:cluster/workspace/:workspace")

	// // ✅ 优先匹配 cluster-scope + name
	// group.Any("/api/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)
	// group.Any("/apis/:group/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)

	// // ✅ namespace 资源 + name
	// group.Any("/api/:version/namespaces/:namespace/:resource/:name", k8sResourceApiGroup.ProxyHandler)
	// group.Any("/apis/:group/:version/namespaces/:namespace/:resource/:name", k8sResourceApiGroup.ProxyHandler)

	// // ✅ 列表请求
	// group.Any("/api/:version/:resource", k8sResourceApiGroup.ProxyHandler)
	// group.Any("/api/:version/namespaces/:namespace/:resource", k8sResourceApiGroup.ProxyHandler)
	// group.Any("/apis/:group/:version/namespaces/:namespace/:resource", k8sResourceApiGroup.ProxyHandler)
	// group.Any("/apis/:group/:version/:resource", k8sResourceApiGroup.ProxyHandler)

	// // 动态获取无组名的资源
	// group.Any("/api/:version/:resource/*name", k8sResourceApiGroup.ProxyHandler)
	// // 无组名namespace级别，例如Pod
	// group.Any("/api/:version/namespaces/:namespace/:resource/*name", k8sResourceApiGroup.ProxyHandler)
	// // 动态获取有组名的资源
	// group.Any("/apis/:group/:version/namespaces/:namespace/:resource/*name", k8sResourceApiGroup.ProxyHandler)
	// // 集群界别资源有组名，无需单独提供接口，后端config.KubeDynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})  可根据提供的namespace为""自动过滤掉，例如：clusterrole
	// group.Any("/apis/:group/:version/:resource/*name", k8sResourceApiGroup.ProxyHandler)

	// // 动态获取无组名的资源
	// group.Any("/api/:version/:resource", k8sResourceApiGroup.ProxyHandler)
	// // 无组名namespace级别，例如Pod
	// group.Any("/api/:version/namespaces/:namespace/:resource", k8sResourceApiGroup.ProxyHandler)
	// // 动态获取有组名的资源
	// group.Any("/apis/:group/:version/namespaces/:namespace/:resource", k8sResourceApiGroup.ProxyHandler)
	// // 集群界别资源有组名，无需单独提供接口，后端config.KubeDynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})  可根据提供的namespace为""自动过滤掉，例如：clusterrole
	// group.Any("/apis/:group/:version/:resource", k8sResourceApiGroup.ProxyHandler)

}
