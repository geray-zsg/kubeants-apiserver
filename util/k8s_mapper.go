package util

import (
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"kubeants.io/models"
)

// GVRKey 用于作为 GVR Map 的键
type GVRKey struct {
	Group    string
	Version  string
	Resource string
}

var GVRToTypeMap = map[GVRKey]reflect.Type{
	{"", "v1", "pods"}:             reflect.TypeOf(v1.Pod{}),
	{"apps", "v1", "deployments"}:  reflect.TypeOf(appsv1.Deployment{}),
	{"apps", "v1", "statefulsets"}: reflect.TypeOf(appsv1.StatefulSet{}),

	{"userbinding.kubeants.io", "v1beta1", "userbindings"}: reflect.TypeOf(models.UserBinding{}),
	// ... 添加更多 CRD 映射
}

// GetStructTypeByGVR 获取 GVR 对应的结构体类型
func GetStructTypeByGVR(group, version, resource string) (reflect.Type, bool) {
	// typ 是对应的 reflect.Type
	// ok 表示该 GVR 是否存在于 GVRToTypeMap 中
	typ, ok := GVRToTypeMap[GVRKey{
		Group:    group,
		Version:  version,
		Resource: resource,
	}]
	return typ, ok
}
