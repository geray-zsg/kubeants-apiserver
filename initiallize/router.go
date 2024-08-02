package initiallize

import (
	"github.com/gin-gonic/gin"
	"kubeant.cn/middleware"
	"kubeant.cn/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors)
	exampleRouterGroup := router.RouterGroupApp.ExampleRouterGroup
	k8sRouterGroup := router.RouterGroupApp.K8SRouterGroup
	exampleRouterGroup.InitExample(r)
	k8sRouterGroup.InitK8SRouter(r)
	return r
}
