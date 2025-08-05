// router/k8s/router.go
package k8s

import (
	"github.com/gin-gonic/gin"
	"kubeants.io/api"
	"kubeants.io/middleware"
)

type K8SRouter struct{}

// func (*K8SRouter) InitK8SRouter(r *gin.Engine) {
// 	k8sResourceApiGroup := api.ApiGroupApp.K8sResourceApi
// 	k8sLogApiGroup := api.ApiGroupApp.LogApi
// 	execApi := api.ApiGroupApp.ExecApi

// 	// ========== Cluster-scope 资源 ==========
// 	cluster := r.Group("/gapi/cluster/:cluster")
// 	cluster.Use(middleware.AuthMiddleware())

// 	{
// 		// 无 Group 的资源，如 Node, PV 等
// 		cluster.Any("/api/:version/:resource", k8sResourceApiGroup.ProxyHandler)
// 		cluster.Any("/api/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)

// 		// 有 Group 的资源，如 CRD
// 		cluster.Any("/apis/:group/:version/:resource", k8sResourceApiGroup.ProxyHandler)
// 		cluster.Any("/apis/:group/:version/:resource/:name", k8sResourceApiGroup.ProxyHandler)
// 	}

// 	// ========== Namespace-scoped 工作空间资源 ==========
// 	workspace := r.Group("/gapi/cluster/:cluster/workspace/:workspace")
// 	workspace.Use(middleware.AuthMiddleware())

// 	{
// 		// 特殊的 /log、/exec 放在前面注册，防止被上面这个拦截（因为 pods/:pod 也匹配 :resource/:name）
// 		// Pod 日志
// 		workspace.GET("/api/v1/namespaces/:namespace/pods/:pod/log", k8sLogApiGroup.GetPodLogs)
// 		// Pod 终端（WebSocket）
// 		workspace.GET("/api/v1/namespaces/:namespace/pods/:pod/exec", execApi.ExecContainer)

// 		// core/v1 命名空间资源（如 Pod、PVC、ConfigMap）
// 		workspace.Any("/api/:version/namespaces/:namespace/:resource", k8sResourceApiGroup.ProxyHandler)
// 		workspace.Any("/api/:version/namespaces/:namespace/:resource/:name", k8sResourceApiGroup.ProxyHandler)

// 		// group API 命名空间资源（如 Deployment、StatefulSet）
// 		workspace.Any("/apis/:group/:version/namespaces/:namespace/:resource", k8sResourceApiGroup.ProxyHandler)
// 		workspace.Any("/apis/:group/:version/namespaces/:namespace/:resource/:name", k8sResourceApiGroup.ProxyHandler)

// 	}
// }

func (*K8SRouter) InitK8SRouter(r *gin.Engine) {

	// 统一使用client-go的ForResource 方案转发k8s原生接口
	k8sResourceApiGroup := api.ApiGroupApp.K8sResourceApi
	k8sLogApiGroup := api.ApiGroupApp.LogApi
	execApi := api.ApiGroupApp.ExecApi
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
	workspace.GET("/api/v1/namespaces/:namespace/pod/:pod/log", k8sLogApiGroup.GetPodLogs)
	// 进入容器终端接口
	workspace.GET("/api/v1/namespaces/:namespace/pod/:pod/exec", execApi.ExecContainer)
}
