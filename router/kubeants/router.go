package kubeants

import (
	"github.com/gin-gonic/gin"
	"kubeants.io/api"
	"kubeants.io/middleware"
)

type KaRouter struct{}

func (*KaRouter) InitKaRouter(r *gin.Engine) {

	group := r.Group("/gapi/cluster/:cluster")
	group.Use(middleware.AuthMiddleware())

	kaApiGroup := api.ApiGroupApp.KaApi

	// 获取用户名下的业务空间列表:/gapi/cluster/:cluster/user/{username}/workspaceslist [get]
	group.GET("/user/:username/workspaceslist", kaApiGroup.WorkspaceListByUser)
	// 获取业务空间下的用户列表
}
