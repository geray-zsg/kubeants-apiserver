package k8s

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"kubeants.io/response"
)

type LogApi struct{}

// GET /gapi/cluster/:cluster/workspace/:workspace/api/v1/namespaces/:namespace/pods/:pod/log
// GET /gapi/.../pods/nginx/log?container=nginx  查看日志
// window.open(`/gapi/.../pods/nginx/log?container=nginx&download=true`)	下载日志
func (*LogApi) GetPodLogs(c *gin.Context) {
	ctx := context.TODO()
	namespace := c.Param("namespace")
	pod := c.Param("pod")

	container := c.Query("container")
	tailStr := c.DefaultQuery("tailLines", "100")
	followStr := c.DefaultQuery("follow", "false")

	tailLines, err := strconv.ParseInt(tailStr, 10, 64)
	if err != nil {
		tailLines = 10000
	}

	follow := followStr == "true"

	logs, err := logService.GetPodLogs(ctx, namespace, pod, container, tailLines, follow)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	// 下载模式
	if c.Query("download") == "true" {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.log", pod))
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, logs)
		return
	}

	// 普通展示
	c.String(200, logs)
}
