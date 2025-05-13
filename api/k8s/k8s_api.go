// /api/k8s/k8s_api.go
package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeants.io/middleware"
	"kubeants.io/response"
	"kubeants.io/util"
	ctrl "sigs.k8s.io/controller-runtime"
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
	ctrl.Log.V(1).Info("这是 debug 日志")
	ctrl.Log.Info("========================> ctrl.Log.Info()...")
	ctrl.Log.V(0).Info("这是 info 日志")

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
		// 定义豁免检查的资源组
		exemptGroups := map[string]bool{
			"workspace.kubeants.io":   true,
			"userbinding.kubeants.io": true,
		}
		// 如果资源属于豁免组，检查请求方法
		if exemptGroups[params.GVR.Group] {
			// 允许 GET 和 LIST 方法跳过权限校验
			if params.K8sVerb != "get" && params.K8sVerb != "list" {
				// 非 GET/LIST 方法，进行权限校验
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
		} else {

			// 非豁免资源，进行权限校验,k8s资源权限校验
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
	setCORSHeaders(c)

	labelSelector := c.Query("labelSelector")
	log.FromContext(ctx).Info("DEBUG selector", "labelSelector", labelSelector)

	if p.IsList {
		list, err := resourceService.ListResources(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, labelSelector)
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		// 删除 managedFields
		for i := range list.Items {
			if metadata, ok := list.Items[i].Object["metadata"].(map[string]interface{}); ok {
				delete(metadata, "managedFields")
			}
		}

		typ, ok := util.GetStructTypeByGVR(p.Group, p.Version, p.Resource)
		if ok {
			listType := reflect.SliceOf(typ)
			listPtr := reflect.New(listType).Interface()

			err = util.UnstructuredListToStructList(list, listPtr)
			if err != nil {
				response.FailWithMessage(c, err.Error())
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"items":      reflect.ValueOf(listPtr).Elem().Interface(),
				"totalItems": reflect.ValueOf(listPtr).Elem().Len(),
			})
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
		// 删除 managedFields
		if metadata, ok := obj.Object["metadata"].(map[string]interface{}); ok {
			delete(metadata, "managedFields")
		}

		// obj 为单个 unstructured.Unstructured 对象
		typ, ok := util.GetStructTypeByGVR(p.Group, p.Version, p.Resource)
		if ok {
			objPtr := reflect.New(typ).Interface()
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, objPtr)
			if err != nil {
				response.FailWithMessage(c, err.Error())
				return
			}

			c.JSON(http.StatusOK, objPtr)
			return
		}

		// fallback: 原始 Unstructured 返回
		c.JSON(http.StatusOK, obj)
	}
}

func handlePost(ctx context.Context, c *gin.Context, p *requestParams, obj *unstructured.Unstructured) {
	setCORSHeaders(c)

	result, err := resourceService.CreateResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, obj)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusCreated, result)
}
func handlePut(ctx context.Context, c *gin.Context, p *requestParams, obj *unstructured.Unstructured) {
	setCORSHeaders(c)

	result, err := resourceService.UpdateResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, obj)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
func handlePatch(ctx context.Context, c *gin.Context, p *requestParams, obj *unstructured.Unstructured) {
	setCORSHeaders(c)

	result, err := resourceService.PatchResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, p.Name, obj)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
func handleDelete(ctx context.Context, c *gin.Context, p *requestParams) {
	setCORSHeaders(c)

	err := resourceService.DeleteResource(ctx, p.Cluster, p.Group, p.Version, p.Resource, p.Namespace, p.Name)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "资源删除成功"})
}

// 设置跨域头部:setCORSHeaders 是 CORS 跨域响应头
// setCORSHeaders 设置标准 CORS 响应头
func setCORSHeaders(c *gin.Context) {
	ctrl.Log.Info("========================> 设置cors响应头...")
	// ctrl.Log.V(1).Info("这是 debug 日志")
	origin := c.Request.Header.Get("Origin")
	if origin != "" {
		fmt.Println("------------------------------------------> 动态代理接口 ProxyHandler 中没有设置 CORS 响应头，统一手动添加跨域响应头...")
		c.Header("Access-Control-Allow-Origin", origin)                                    // 允许来源
		c.Header("Access-Control-Allow-Credentials", "true")                               // 支持 Cookie、Token
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Token")   // 支持的请求头
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS") // 支持的方法
		c.Header("Access-Control-Expose-Headers", "Content-Length, New-Token")             // 客户端可见的响应头
		fmt.Println("------------------------------------------> 动态代理接口 ProxyHandler 中没有设置 CORS 响应头，统一手动添加跨域响应头结束")
		ctrl.Log.Info("========================> 设置跨域响应头成功！！！！！！！！！！！！！！！！！！！！！！！！")
	}
	fmt.Printf("------------------------------------------> 动态代理接口 ProxyHandler 中设置 CORS 响应头信息 c.Request.Header.Get === Origin：%v \n", origin)
	ctrl.Log.Info("========================> 设置跨域响应头成功", "origin", origin)
}
