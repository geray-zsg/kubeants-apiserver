package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"kubeants.io/config"
)

func Cors(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")
	trimmedOrigin := strings.TrimSuffix(origin, "/")

	allowedOrigins := config.CONF.Cors.AllowedOrigins
	defaultOrigin := config.CONF.Cors.DefaultOrigins
	isAllowed := false

	for _, allowed := range allowedOrigins {
		if trimmedOrigin == strings.TrimSuffix(allowed, "/") {
			isAllowed = true
			break
		}
	}

	// 设置跨域响应头
	allowedOrigin := defaultOrigin
	if isAllowed {
		allowedOrigin = origin
	}
	setCorsHeaders(c, allowedOrigin)

	// OPTIONS 预检请求：立即返回
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}

func setCorsHeaders(c *gin.Context, origin string) {
	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Headers", config.CONF.Cors.AccessControlAllowHeaders)
	c.Header("Access-Control-Allow-Methods", config.CONF.Cors.AccessControlAllowMethods)
	c.Header("Access-Control-Expose-Headers", config.CONF.Cors.AccessControlExposeHeaders)
	c.Header("Access-Control-Allow-Credentials", config.CONF.Cors.AccessControlAllowCredentials)
}
