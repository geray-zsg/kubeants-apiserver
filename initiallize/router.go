package initiallize

import (
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kubeants.io/config"
	"kubeants.io/docs"
	"kubeants.io/middleware"
	"kubeants.io/router"
	ctrl "sigs.k8s.io/controller-runtime"
)

func Routers() *gin.Engine {
	// r := gin.Default()
	r := gin.New()
	r.RedirectTrailingSlash = false // 禁用自动斜杠重定向
	r.Use(gin.Logger(), gin.Recovery())

	corsEable := config.CONF.Cors.Enable
	ctrl.Log.V(0).Info("使用使用配置文件中的跨域配置", "是否启用跨域配置", corsEable)
	if corsEable {
		r.Use(middleware.Cors)
	} else {
		r.Use(cors.Default())
	}

	// r.Use(gin.LoggerWithFormatter(CustomLogFormatter))
	group := r.Group("/gapi")
	{
		// 初始化 Swagger 信息
		docs.SwaggerInfo.Title = "蚁挚"
		docs.SwaggerInfo.Description = "kubeants"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.BasePath = "/gapi"

		group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		group.POST("/system/login", middleware.LoginHandler)

	}

	userGroup := group.Group("/system")
	userGroup.Use(middleware.AuthMiddleware())
	userGroup.GET("/userinfo/:username", middleware.GetUserInfo)

	exampleRouterGroup := router.RouterGroupApp.ExampleRouterGroup
	k8sRouterGroup := router.RouterGroupApp.K8SRouterGroup
	kaRouterGroup := router.RouterGroupApp.KaRouterGroup
	exampleRouterGroup.InitExample(r)
	k8sRouterGroup.InitK8SRouter(r)
	kaRouterGroup.InitKaRouter(r)

	return r
}

// 自定义日志格式
func CustomLogFormatter(param gin.LogFormatterParams) string {
	// 格式化日志：时间戳、HTTP方法、路径、状态码、响应时间、客户端IP
	return "[" + param.TimeStamp.Format(time.RFC3339) + "] " +
		param.Method + " " +
		param.Path + " " +
		strconv.Itoa(param.StatusCode) + " " +
		param.ClientIP + " " +
		param.Latency.String()
}
