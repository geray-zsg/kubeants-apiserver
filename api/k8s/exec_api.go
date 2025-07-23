package k8s

import (
	"github.com/gin-gonic/gin"
	"kubeants.io/response"
	"kubeants.io/service/k8s"
)

type ExecApi struct{}

// GET /gapi/cluster/:cluster/workspace/:workspace/api/v1/namespaces/:namespace/pods/:pod/exec
func (*ExecApi) ExecContainer(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("pod")
	container := c.DefaultQuery("container", "")
	command := c.DefaultQuery("command", "/bin/sh")

	wsConn, err := k8s.UpgradeToWebSocket(c.Writer, c.Request)
	if err != nil {
		response.FailWithMessage(c, "升级 WebSocket 失败: "+err.Error())
		return
	}
	defer wsConn.Close()

	err = k8s.NewExecService().ExecToPod(wsConn, namespace, podName, container, command)
	if err != nil {
		wsConn.WriteMessage(1, []byte("连接失败："+err.Error()))
	}
}
