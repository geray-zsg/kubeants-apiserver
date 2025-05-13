package util

import (
	"encoding/json"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// 辅助方法：将 Unstructured 转换为目标结构体
func UnstructuredToStruct(u *unstructured.Unstructured, out interface{}) error {
	data, err := u.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

// 辅助函数：转换Unstructured list为目标结构列表
func UnstructuredListToStructList(u *unstructured.UnstructuredList, out interface{}) error {
	// 检查输出是否为切片指针
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr || outVal.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("output must be a pointer to a slice")
	}

	// 获取切片元素类型
	elemType := outVal.Elem().Type().Elem()

	// 创建结果切片
	result := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(u.Items))

	// 遍历UnstructuredList中的每个项目
	for _, item := range u.Items {
		// 为每个项目创建新的结构体实例
		elem := reflect.New(elemType).Interface()

		// 将Unstructured转换为结构体
		err := UnstructuredToStruct(&item, elem)
		if err != nil {
			return err
		}

		// 将转换后的结构体添加到结果切片中
		result = reflect.Append(result, reflect.ValueOf(elem).Elem())
	}

	// 将结果设置到输出参数中
	outVal.Elem().Set(result)
	return nil
}
