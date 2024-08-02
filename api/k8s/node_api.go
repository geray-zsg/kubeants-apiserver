package k8s

import (
	"context"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeant.cn/config"
	"kubeant.cn/response"
)

type NodeApi struct{}

func (n *NodeApi) GetNodeListOrDetail(c *gin.Context) {
	ctx := context.TODO()
	name := c.Param("name")
	if name != "" {
		// Node详情
		nodeDetail, err := n.GetNodeDetail(ctx, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "Node信息查询成功", nodeDetail)
		return
	}

	// Node列表
	nodeList, err := config.KubeClientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	items := nodeList.Items
	totailItems := len(items)
	resp := gin.H{
		"items":       items,
		"totailItems": totailItems,
	}

	response.SuccessWithDetailed(c, "Node列表信息查询成功", resp)

}

func (*NodeApi) GetNodeDetail(ctx context.Context, name string) (*corev1.Node, error) {
	nodeDetail, err := config.KubeClientSet.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return nodeDetail, nil
}
