package k8s

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeants.io/config"
	"kubeants.io/response"
	"kubeants.io/util"
)

type PodApi struct{}

// Get all pods in all namespace
func (*PodApi) GetAllPods(c *gin.Context) {
	ctx := context.TODO()

	podList, err := config.KubeClientSet.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	items := podList.Items
	totalItems := len(items)
	resp := gin.H{
		"items":      items,
		"totalItems": totalItems,
	}

	c.JSON(http.StatusOK, resp)
}

// Get all pods in  namespace detail
func (p *PodApi) GetPodsInNamespaceORDerail(c *gin.Context) {
	ctx := context.TODO()

	namespace := c.Param("namespace") //从路径获取Namespace（必填项否则包404）
	name := c.Param("name")

	if name != "" {
		// 查看特定 Pod 详情
		// oldPodDetail, err := p.GetoldPodDetail(ctx, namespace, name)
		oldPodDetail, err := podService.GetoldPodDetail(ctx, namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		response.SuccessWithDetailed(c, "Pod信息查询成功", oldPodDetail)
		return
	}

	// 否则查看namespace下的Pods list
	podList, err := config.KubeClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	items := podList.Items
	totalItems := len(items)
	resp := gin.H{
		"items":      items,
		"totalItems": totalItems,
	}
	response.SuccessWithDetailed(c, "所有Pod信息查询成功", resp)
}

// Create Or Update Pod
func (p *PodApi) CreateOrUpdatePod(c *gin.Context) {
	ctx := context.TODO()
	namespace := c.Param("namespace")

	dataByteBodyPod, err := util.RequestUnmarshalForJSONORYAML[corev1.Pod](c)
	if err != nil {
		response.FailWithMessage(c, err.Error())
	}

	// 设置命名空间
	if dataByteBodyPod.Namespace == "" {
		if namespace != "" {
			dataByteBodyPod.Namespace = namespace
		} else {
			dataByteBodyPod.Namespace = "default"
		}
	}

	// 最好有什么方法去校验一下，如果没有任何改变则跳过不操作
	// updatePod, err := p.deleteAndCreatePod(ctx, &dataByteBodyPod)
	updatePod, err := podService.DeleteAndCreatePod(ctx, dataByteBodyPod)
	if err != nil {
		response.FailWithMessage(c, fmt.Sprintf("Failed to update pod: %v", err))
	}

	response.SuccessWithDetailed(c, "Pod创建或修改成功", updatePod)
}

// Delete pod
func (*PodApi) DeletePod(c *gin.Context) {
	ctx := context.TODO()
	namespace := c.Param("namespace")
	name := c.Param("name")

	// 传入错误的namespace和name会找不到所以不需要判断

	err := podService.DeletePod(ctx, namespace, name)
	if err != nil {
		response.FailWithMessage(c, "Failed to delete pod: "+err.Error())
		return
	}
	response.Success(c)
}
