package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"kubeants.com/response"
)

// ResourceApi 处理 Kubernetes 资源的 API
type ResourceApi struct{}

// ProxyHandler 代理所有 HTTP 方法
func (*ResourceApi) ProxyHandler(c *gin.Context) {
	ctx := context.TODO()
	cluster := c.Param("cluster")
	group := c.Param("group")
	version := c.Param("version")
	resource := c.Param("resource")
	namespace := c.Param("namespace")
	name := c.Param("name")
	// 使用 strings.Trim 函数去除开头和结尾的 '/'
	name = strings.Trim(name, "/")
	fmt.Println("Param name:", name)
	if name != "" && !isValidKubernetesName(name) {
		response.FailWithMessage(c, "name名称不符合k8s命名规范：1.名称只能包含小写字母、数字、连字符（-）和点（.）;2.名称必须以字母或数字开头和结尾;3.名称中的连字符（-）不能连续出现，且不能位于名称的开头或结尾;4.名称的长度必须在 1 到 63 个字符之间;5.名称不能以点（.）结尾。"+name)
		return
	}

	// 解析请求体（如果有）
	var obj unstructured.Unstructured
	if c.Request.Body != nil && c.Request.ContentLength > 0 {
		defer c.Request.Body.Close()

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			response.FailWithMessage(c, "Failed to read request body")
			return
		}

		if err := json.Unmarshal(body, &obj); err != nil {
			response.FailWithMessage(c, fmt.Sprintf("Invalid JSON format: %v", err))
			return
		}
	}

	// 识别 HTTP 方法并调用 Service 层
	switch c.Request.Method {
	case http.MethodGet:
		if name == "" {
			resources, err := resourceService.ListResources(ctx, cluster, group, version, resource, namespace)
			if err != nil {
				response.FailWithMessage(c, err.Error())
				return
			}
			c.JSON(http.StatusOK, resources)
		} else {
			obj, err := resourceService.GetResource(ctx, cluster, group, version, resource, namespace, name)
			if err != nil {
				response.FailWithMessage(c, err.Error())
				return
			}
			c.JSON(http.StatusOK, obj)
		}
	case http.MethodPost:
		result, err := resourceService.CreateResource(ctx, cluster, group, version, resource, namespace, &obj)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		c.JSON(http.StatusCreated, result)
	case http.MethodPut:
		result, err := resourceService.UpdateResource(ctx, cluster, group, version, resource, namespace, &obj)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	case http.MethodPatch:
		result, err := resourceService.PatchResource(ctx, cluster, group, version, resource, namespace, name, &obj)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	case http.MethodDelete:
		err := resourceService.DeleteResource(ctx, cluster, group, version, resource, namespace, name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Resource deleted"})
	case http.MethodOptions:
		// 返回支持的方法
		c.Header("Allow", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Status(http.StatusNoContent)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": fmt.Sprintf("Method %s not supported", c.Request.Method)})
	}
}

// GetResourcesHandler 通用资源查询
func (*ResourceApi) GetResourcesHandler(c *gin.Context) {
	ctx := context.TODO()
	group := c.Param("group")
	version := c.Param("version")
	resource := c.Param("resource")

	resources, err := resourceService.Resources(ctx, group, version, resource)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resources,
	})
}

// isValidKubernetesName 校验一个名称是否符合 Kubernetes 的命名规则
func isValidKubernetesName(name string) bool {
	// Kubernetes 名称的正则表达式
	// 注意：这个正则表达式可能需要根据 Kubernetes 的实际规则进行调整
	var validName = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	// 长度检查
	if len(name) < 1 || len(name) > 63 {
		return false
	}
	// 校验名称是否匹配正则表达式
	return validName.MatchString(name)
}
