// /service/k8s/k8s_service.go
package k8s

import (
	"context"
	"encoding/json"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"kubeants.io/config"
)

type ResourceService struct{}

// 这里使用 dynamicClient.Resource(gvr).Namespace(v1.NamespaceAll).List(...) 来获取所有资源。
// 这样就不用为每个资源单独写 PodService、DeploymentService 之类的代码。
// GetResources(已废弃Resources) 动态获取 K8s 资源
func (*ResourceService) Resources(ctx context.Context, group, version, resource string) ([]unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	// 获取动态客户端
	dynamicClient := config.KubeDynamicClient

	// 查询所有命名空间资源
	resourceList, err := dynamicClient.Resource(gvr).Namespace(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return resourceList.Items, nil
}

// ListResources 获取资源列表
func (s *ResourceService) ListResources(ctx context.Context, cluster, group, version, resource, namespace, labelSelector string) (*unstructured.UnstructuredList, error) {

	gvr := getGVR(group, version, resource)
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}

	return config.KubeDynamicClient.Resource(gvr).Namespace(namespace).List(ctx, listOptions)
}

// ListResourcesByLabelSelector 通过Label选择器表达式获取资源列表
func (s *ResourceService) ListResourcesByLabelSelector(ctx context.Context, cluster, group, version, resource, namespace, labelSelector string) (*unstructured.UnstructuredList, error) {
	gvr := getGVR(group, version, resource)

	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	return config.KubeDynamicClient.Resource(gvr).Namespace(namespace).List(ctx, listOptions)
}

// GetResource 获取单个资源
func (s *ResourceService) GetResource(ctx context.Context, cluster, group, version, resource, namespace, name string) (*unstructured.Unstructured, error) {
	gvr := getGVR(group, version, resource)
	return config.KubeDynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateResource 创建资源
func (s *ResourceService) CreateResource(ctx context.Context, cluster, group, version, resource, namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {

	gvr := getGVR(group, version, resource)
	return config.KubeDynamicClient.Resource(gvr).Namespace(namespace).Create(ctx, obj, metav1.CreateOptions{})
}

// UpdateResource 更新资源（PUT 方法）
func (s *ResourceService) UpdateResource(ctx context.Context, cluster, group, version, resource, namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {

	gvr := getGVR(group, version, resource)
	return config.KubeDynamicClient.Resource(gvr).Namespace(namespace).Update(ctx, obj, metav1.UpdateOptions{})
}

// PatchResource 部分更新资源（PATCH 方法）
func (s *ResourceService) PatchResource(ctx context.Context, cluster, group, version, resource, namespace, name string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {

	gvr := getGVR(group, version, resource)
	// 获取资源名称
	if name == "" {
		name = obj.GetName()
	}

	if name == "" {
		return nil, errors.New("resource name is required for patching")
	}

	// 将 obj.Object 转换为 JSON 格式的 []byte
	patchData, err := json.Marshal(obj.Object)
	if err != nil {
		return nil, err
	}

	// 选择 Patch 类型，通常是 MergePatch 或 StrategicMergePatch
	patchType := types.MergePatchType // 也可以用 types.StrategicMergePatchType
	return config.KubeDynamicClient.Resource(gvr).Namespace(namespace).Patch(ctx, name, patchType, patchData, metav1.PatchOptions{})
}

// DeleteResource 删除资源
func (s *ResourceService) DeleteResource(ctx context.Context, cluster, group, version, resource, namespace, name string) error {

	gvr := getGVR(group, version, resource)
	return config.KubeDynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// getGVR 构造 GVR（GroupVersionResource）
func getGVR(group, version, resource string) (gvr schema.GroupVersionResource) {
	if group == "" {
		gvr = schema.GroupVersionResource{
			Group:    "", // 无组名资源没有 Group 部分
			Version:  version,
			Resource: resource,
		}
		return gvr
	}

	gvr = schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	return gvr
}
