package util

import (
	"encoding/json"
	"fmt"

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
