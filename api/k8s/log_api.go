package k8s

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"kubeants.io/response"
)

type LogApi struct{}

// GET /gapi/cluster/:cluster/workspace/:workspace/api/v1/namespaces/:namespace/pods/:pod/log
func (*LogApi) GetPodLogs(c *gin.Context) {
	ctx := context.TODO()
	namespace := c.Param("namespace")
	pod := c.Param("pod")

	// 可选参数
	container := c.Query("container") // 支持多容器指定
	tailStr := c.DefaultQuery("tailLines", "100")
	followStr := c.DefaultQuery("follow", "false")

	tailLines, err := strconv.ParseInt(tailStr, 10, 64)
	if err != nil {
		tailLines = 100
	}

	follow := followStr == "true"

	// logService := k8s.LogService{}
	logs, err := logService.GetPodLogs(ctx, namespace, pod, container, tailLines, follow)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	c.String(200, logs)
}
