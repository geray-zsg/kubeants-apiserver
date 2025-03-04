package initiallize

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"kubeants.io/middleware"
	"kubeants.io/router"
)

func Routers() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors)
	// r.Use(gin.LoggerWithFormatter(CustomLogFormatter))

	exampleRouterGroup := router.RouterGroupApp.ExampleRouterGroup
	k8sRouterGroup := router.RouterGroupApp.K8SRouterGroup
	exampleRouterGroup.InitExample(r)
	k8sRouterGroup.InitK8SRouter(r)

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
