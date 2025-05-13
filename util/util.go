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

// // 辅助方法：将 Unstructured 转换为目标结构体
// func UnstructuredToStruct(u *unstructured.Unstructured, out interface{}) error {
// 	data, err := u.MarshalJSON()
// 	if err != nil {
// 		return err
// 	}
// 	return json.Unmarshal(data, out)
// }

// // 辅助函数：转换Unstructured list为目标结构列表
// func UnstructuredListToStructList(u *unstructured.UnstructuredList, out interface{}) error {
// 	// 检查输出是否为切片指针
// 	outVal := reflect.ValueOf(out)
// 	if outVal.Kind() != reflect.Ptr || outVal.Elem().Kind() != reflect.Slice {
// 		return fmt.Errorf("output must be a pointer to a slice")
// 	}

// 	// 获取切片元素类型
// 	elemType := outVal.Elem().Type().Elem()

// 	// 创建结果切片
// 	result := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(u.Items))

// 	// 遍历UnstructuredList中的每个项目
// 	for _, item := range u.Items {
// 		// 为每个项目创建新的结构体实例
// 		elem := reflect.New(elemType).Interface()

// 		// 将Unstructured转换为结构体
// 		err := UnstructuredToStruct(&item, elem)
// 		if err != nil {
// 			return err
// 		}

// 		// 将转换后的结构体添加到结果切片中
// 		result = reflect.Append(result, reflect.ValueOf(elem).Elem())
// 	}

// 	// 将结果设置到输出参数中
// 	outVal.Elem().Set(result)
// 	return nil
// }
