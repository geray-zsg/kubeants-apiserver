// /api/k8s/k8s_api.go
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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeants.io/middleware"
	"kubeants.io/response"
	"kubeants.io/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ResourceApi 处理 Kubernetes 资源的 API
type ResourceApi struct{}

// 定义请求参数
type requestParams struct {
	Cluster     string
	Workspace   string
	Group       string
	Version     string
	Resource    string
	Namespace   string
	Name        string
	GVR         schema.GroupVersionResource
	K8sVerb     string
	IsList      bool
	Username    string
	SAToken     string
	RequestBody []byte
}

// ProxyHandler 代理所有 HTTP 方法
func (*ResourceApi) ProxyHandler(c *gin.Context) {
	ctx := context.TODO()
	logger := log.FromContext(ctx)

	// 1. 解析请求参数
	params, err := parseRequestParams(ctx, c)
	if err != nil {
		logger.Error(err, "参数解析失败")
		response.FailWithMessage(c, err.Error())
		return
	}
	logger.Info("接收到请求", "params", params)

	// 2. 非 admin 用户权限校验
	if params.Username != "admin" {
		allowed, err := middleware.CheckResourcePermission(ctx, params.SAToken, params.GVR, params.K8sVerb, params.Namespace)
		if err != nil {
			logger.Error(err, "权限校验出错")
			response.FailWithMessage(c, "权限校验失败")
			return
		}
		if !allowed {
			logger.Info("权限不足，当前用户无访问该资源权限")
			response.FailWithMessage(c, "权限不足，当前用户无访问该资源权限")
			return
		}
	}

	// 3. 如果需要 body，解析 body
	var obj unstructured.Unstructured
	if c.Request.ContentLength > 0 {
		if err := parseRequestBody(c, &obj); err != nil {
			logger.Error(err, "请求体解析失败")
			response.FailWithMessage(c, err.Error())
			return
		}
	}

	// 4. 分发不同方法处理
	switch c.Request.Method {
	case http.MethodGet:
		handleGet(ctx, c, params)
	case http.MethodPost:
		handlePost(ctx, c, params, &obj)
	case http.MethodPut:
		handlePut(ctx, c, params, &obj)
	case http.MethodPatch:
		handlePatch(ctx, c, params, &obj)
	case http.MethodDelete:
		handleDelete(ctx, c, params)
	case http.MethodOptions:
		c.Header("Allow", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Status(http.StatusNoContent)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": fmt.Sprintf("Method %s not supported", c.Request.Method)})
	}
}

// ========== 参数解析 ==========
func parseRequestParams(ctx context.Context, c *gin.Context) (*requestParams, error) {
	logger := log.FromContext(ctx)
	name := strings.Trim(c.Param("name"), "/")
	isList := name == ""
	// 转换http请求为k8s的动词
	verb := util.HTTPMethodToK8sVerb(c.Request.Method, isList, c.Request.URL.Query())

	saToken, exists := c.Get("sa_token")
	if !exists {
		logger.Info("未获取到用户身份认证信息（sa_token）")
		return nil, fmt.Errorf("未获取到用户身份认证信息（sa_token）")
	}

	username, exists := c.Get("username")
	if !exists {
		logger.Info("未获取到用户名信息")
		return nil, fmt.Errorf("未获取到用户名信息")
	}

	if name != "" && !isValidKubernetesName(name) {
		logger.Info("资源名称不符合 Kubernetes 命名规范")
		return nil, fmt.Errorf("资源名称不符合 Kubernetes 命名规范: %s", name)
	}

	group := c.Param("group")
	version := c.Param("version")
	resource := c.Param("resource")

	logger.Info("请求详情",
		"cluster", c.Param("cluster"),
		"workspace", c.Param("workspace"),
		"group", group,
		"version", version,
		"resource", resource,
		"namespace", c.Param("namespace"),
		"name", name,
		"query", c.Request.URL.RawQuery,
	)

	return &requestParams{
		Cluster:   c.Param("cluster"),
		Workspace: c.Param("workspace"),
		Group:     group,
		Version:   version,
		Resource:  resource,
		Namespace: c.Param("namespace"),
		Name:      name,
		IsList:    isList,
		K8sVerb:   verb,
		Username:  username.(string),
		SAToken:   saToken.(string),
		GVR: schema.GroupVersionResource{
			Group:    group,
			Version:  version,
			Resource: resource,
		},
	}, nil

}

// isValidKubernetesName 校验一个名称是否符合 Kubernetes 的命名规则
func isValidKubernetesName(name string) bool {
	// 长度检查
	if len(name) < 1 || len(name) > 63 {
		return false
	}
	// Kubernetes 名称的正则表达式
	// 注意：这个正则表达式可能需要根据 Kubernetes 的实际规则进行调整
	var validName = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

	// 校验名称是否匹配正则表达式
	return validName.MatchString(name)
}

func parseRequestBody(c *gin.Context, obj *unstructured.Unstructured) error {
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("读取请求体失败: %v", err)
	}
	if err := json.Unmarshal(body, obj); err != nil {
		return fmt.Errorf("请求体格式错误: %v", err)
	}
	return nil
}

// ========== 不同方法的处理逻辑 ==========
func handleGet(ctx context.Context, c *gin.Context, p *requestParams) {
	labelSelector := c.Query("labelSelector")
	log.FromContext(ctx).Info("DEBUG selector", "labelSelector", labelSelector)

	if p.IsList {
		list, err := resourceService.ListResources(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, labelSelector)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		resp := gin.H{
			"items":      list,
			"totalItems": len(list.Items),
		}
		c.JSON(http.StatusOK, resp)
		// c.JSON(http.StatusOK, list)
	} else {
		obj, err := resourceService.GetResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, p.Name)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, obj)
	}
}

func handlePost(ctx context.Context, c *gin.Context, p *requestParams, obj *unstructured.Unstructured) {
	result, err := resourceService.CreateResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, obj)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusCreated, result)
}
func handlePut(ctx context.Context, c *gin.Context, p *requestParams, obj *unstructured.Unstructured) {
	result, err := resourceService.UpdateResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, obj)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
func handlePatch(ctx context.Context, c *gin.Context, p *requestParams, obj *unstructured.Unstructured) {
	result, err := resourceService.PatchResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, p.Name, obj)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
func handleDelete(ctx context.Context, c *gin.Context, p *requestParams) {
	err := resourceService.DeleteResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, p.Name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "资源删除成功"})
}
