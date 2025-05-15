package models

import "k8s.io/apimachinery/pkg/runtime/schema"

// 定义请求参数
type RequestParams struct {
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
