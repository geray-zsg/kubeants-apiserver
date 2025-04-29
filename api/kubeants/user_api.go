package kubeants

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"kubeants.io/models"
	"kubeants.io/response"
	"kubeants.io/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type UserApi struct{}

// GetWorkspaceListByUser 获取user的workspace列表
// @Summary 认证通过后返回用户信息。
// @Description 验证用户输入的用户名和密码或JWTtoekn，如果正确，则返回数据。
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "用户的workspace列表获取成功"
// @Router /gapi/cluster/{cluster}/user/{username}/workspaceslist [get]
func (u *UserApi) WorkspaceListByUser(c *gin.Context) {
	ctx := context.TODO()
	logger := log.FromContext(ctx)
	username := c.Param("username")

	logger.Info("GetWorkspaceListByUser", "username", username)

	bindingList, err := u.getUserBindingByLabel(ctx, username)
	if err != nil {
		logger.Error(err, "userbindingList 获取失败")
		response.FailWithMessage(c, "userbindingList 获取失败"+err.Error())
		return
	}

	msg := "userbindingList 数据获取成功"
	resp := gin.H{
		"msg":   msg,
		"data":  bindingList,
		"tatol": len(bindingList),
	}

	logger.Info(msg)
	c.JSON(http.StatusOK, resp)
}

// 获取userbinding通过label
func (*UserApi) getUserBindingByLabel(ctx context.Context, username string) (bindinglist []models.UserBinding, err error) {
	group := "userbinding.kubeants.io"
	version := "v1beta1"
	resource := "userbindings"

	labelSelector := "kubeants.io/user=" + username

	list, err := resourceService.ListResourcesByLabelSelector(ctx, "", group, version, resource, "", labelSelector)
	if err != nil {
		return nil, err
	}

	// 将 userunstructured 转换为 UserBinding 结构体
	if err = util.UnstructuredListToStructList(list, &bindinglist); err != nil {
		return nil, err
	}
	return bindinglist, nil
}
