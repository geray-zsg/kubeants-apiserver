package k8s

import (
	"context"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeants.io/config"
	"kubeants.io/response"
	"kubeants.io/util"
)

type NamespaceApi struct{}

// Get namespace list or namespace detail
func (n *NamespaceApi) GetNamespaceList(c *gin.Context) {

	ctx := context.TODO()
	name := c.Param("name")
	if name != "" {
		nsGet, err := namespaceService.GetNamespaceDetail(ctx, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "namespace信息查询成功", nsGet)
		return
	}

	namespaceList, err := config.KubeClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	items := namespaceList.Items
	totalItems := len(items)
	resp := gin.H{
		"items":      items,
		"totalItems": totalItems,
	}

	response.SuccessWithDetailed(c, "namespace列表信息查询成功", resp)
}

func (*NamespaceApi) DeleteNamespace(c *gin.Context) {
	ctx := context.TODO()
	name := c.Param("name")
	err := namespaceService.DeleteNamespace(ctx, name)
	if err != nil {
		response.FailWithMessage(c, "Failed to delete namespace:"+err.Error())
	}
	response.Success(c)
}

func (*NamespaceApi) CreateNamespace(c *gin.Context) {
	ctx := context.TODO()
	dataByteBody, err := util.RequestUnmarshalForJSONORYAML[corev1.Namespace](c)
	if err != nil {
		response.FailWithMessage(c, "Failed to unmarshal request body:"+err.Error())
		return
	}
	createNamespace, err := namespaceService.CreateNamespace(ctx, dataByteBody)
	if err != nil {
		response.FailWithMessage(c, "Failed to create namespace:"+err.Error())
		return
	}
	response.SuccessWithDetailed(c, "namespace创建成功", createNamespace)
}
