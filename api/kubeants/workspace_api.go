package kubeants

import (
	"context"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type WorkspaceApi struct{}

// Workspace 获取 Workspace 信息
// @Summary 获取 Workspace 信息
// @Description 认证通过后获取 Workspace 信息（不传 name 时返回列表，传 name 时返回详情）
// @Tags workspace
// @Accept json
// @Produce json
// @Param cluster path string true "集群名称"
// @Param name path string false "Workspace 名称（可选）"
// @Success 200 {object} map[string]interface{} "登陆成功"
// @Router /gapi/cluster/{cluster}/workspace [get]
// @Router /gapi/cluster/{cluster}/workspace/{name} [get]
func (*UserApi) Workspace(c *gin.Context) {
	ctx := context.TODO()
	logger := log.FromContext(ctx)
	cluster := c.Param("cluster")
	name := c.Param("name") // 如果是路径参数

	logger.Info("GetUserListByWorkspace", "cluster", cluster, "workspace", name)

	if name == "" {
		// 获取所有 Workspace 列表

	} else {
		// 获取单个 Workspace 详情
		// getWorkspaceDetail(c, cluster, name)
	}

}

// GetUserListByWorkspace 获取workspace下user列表
// @Summary 认证通过后返回用户信息。
// @Description 验证用户输入的用户名和密码，如果正确，则返回 JWT。
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "登陆成功"
// @Router /gapi/cluster/{cluster}/workspace/{workspace}/userlist [get]
func (*UserApi) GetUserListByWorkspace(c *gin.Context) {
	ctx := context.TODO()
	logger := log.FromContext(ctx)
	workspace := c.Param("workspace")

	logger.Info("GetUserListByWorkspace", "workspace", workspace)

}
