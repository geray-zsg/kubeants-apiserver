package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kubeants.io/config"
)

func Cors(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")

	// allowedOrigins := []string{"http://demo-mysql.10.179.25.1.nip.io", "http://localhost:8080"}
	allowedOrigins := config.CONF.Cors.AllowedOrigins
	fmt.Printf("Origin: %s, allowedOrigins: %v", origin, allowedOrigins)
	// 校验 Origin
	isAllowed := false
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			isAllowed = true
			break
		}
	}
	if isAllowed {
		c.Header("Access-Control-Allow-Origin", origin)
	} else {
		c.Header("Access-Control-Allow-Origin", config.CONF.Cors.DefaultOrigins) // 允许所有域名，不建议用于生产
		// c.Header("Access-Control-Allow-Origin", "http://localhost:8080") // 默认值
	}

	c.Header("Access-Control-Allow-Headers", config.CONF.Cors.AccessControlAllowHeaders)
	c.Header("Access-Control-Allow-Methods", config.CONF.Cors.AccessControlAllowMethods)
	c.Header("Access-Control-Expose-Headers", config.CONF.Cors.AccessControlExposeHeaders)
	c.Header("Access-Control-Allow-Credentials", config.CONF.Cors.AccessControlAllowCredentials)
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}

	// 打印头header头信息
	for k, v := range c.Request.Header {
		fmt.Printf("Header: %s: %v\n", k, v)
	}
	// 处理请求
	c.Next()
}
