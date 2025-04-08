package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func RequestUnmarshalForJSONORYAML[T any](c *gin.Context) (*T, error) {

	// 读取请求体
	dataByteBody, err := c.GetRawData()
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	// 尝试将请求体解析为JSON或者YAML格式的Pod对象
	// var dataByteBodyPod corev1.Pod
	// 创建目标对象
	var obj T

	if err := json.Unmarshal(dataByteBody, &obj); err != nil {
		// 如果JSON解析失败，尝试用YML解析
		if err := yaml.Unmarshal(dataByteBody, &obj); err != nil {
			return nil, fmt.Errorf("invalid JSON or YAML format: %v", err)
		}
	}

	return &obj, nil
}

// HTTPMethodToK8sVerb 将 HTTP 方法映射为 Kubernetes 的权限 verb。
func HTTPMethodToK8sVerb(method string, isResourceList bool, queryParams map[string][]string) string {
	switch method {
	case http.MethodGet:
		// watch 通常通过 GET 请求 + ?watch=true 实现
		if watchParams, ok := queryParams["watch"]; ok && len(watchParams) > 0 && (watchParams[0] == "true" || watchParams[0] == "1") {
			return "watch"
		}
		if isResourceList {
			return "list"
		}
		return "get"

	case http.MethodPost:
		// proxy/create/eviction/exec/attach 都是 POST 请求，根据 URL 和资源判断，此处统一为 create
		return "create"

	case http.MethodPut:
		return "update"

	case http.MethodPatch:
		return "patch"

	case http.MethodDelete:
		// deletecollection 属于 DELETE 方法 + 无 name（即删除列表）
		if isResourceList {
			return "deletecollection"
		}
		return "delete"

	case http.MethodConnect:
		// 用于一些特殊场景，比如 `exec`, `port-forward`, `attach` 等
		return "connect"
	}

	return ""
}
