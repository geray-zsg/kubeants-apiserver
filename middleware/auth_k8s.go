// /middleware/auth_k8s.go
package middleware

import (
	"context"
	"fmt"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kubeants.io/config"
	"kubeants.io/models"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CheckResourcePermission 校验 Token 是否具备指定资源的操作权限
func CheckResourcePermission(ctx context.Context, saToken string, gvr schema.GroupVersionResource, verb, namespace string) (bool, error) {
	logger := log.FromContext(ctx)
	// 1. 构建 SelfSubjectAccessReview 请求
	logger.Info("构建用户的SelfSubjectAccessReview，准备验证权限,方法收到的请求有", "gvr", gvr, "verb", verb, "namespace", namespace)
	ssar := &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Group:     gvr.Group,
				Version:   gvr.Version,
				Resource:  gvr.Resource,
				Verb:      verb,
				Namespace: namespace, // 若为空则为集群范围权限
			},
		},
	}

	logger.Info("构建用户的SelfSubjectAccessReview，准备验证权限", "SelfSubjectAccessReview", ssar)

	// 2. 调用 API 校验权限
	tokenClientSet, err := NewClientWithSAToken(saToken)
	if err != nil {
		return false, fmt.Errorf("无法创建 Token 客户端: %v", err)
	}
	logger.Info("创建 Token 客户端成功", "tokenClientSet", tokenClientSet)

	response, err := tokenClientSet.AuthorizationV1().SelfSubjectAccessReviews().Create(context.TODO(), ssar, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "无法创建 SelfSubjectAccessReview 请求")
		return false, fmt.Errorf("权限校验请求失败: %v", err)
	}

	// 3. 返回权限状态
	logger.Info("权限校验结果", "response", response.Status.Allowed)
	return response.Status.Allowed, nil
}

// 用 saToken 生成 clientSet
func NewClientWithSAToken(saToken string) (*kubernetes.Clientset, error) {
	config := &rest.Config{
		Host:        config.Kubeconfig, // 这里从已有的配置中获取 API Server 地址
		BearerToken: saToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // 若使用自签名证书，可设为 true
		},
	}
	return kubernetes.NewForConfig(config)
}

// 判断是否为豁免资源的函数
func IsExemptedResource(params *models.RequestParams) bool {
	for _, r := range config.CONF.Authz.ExemptResources {
		if (r.Group == "*" || r.Group == params.GVR.Group) &&
			(r.Version == "" || r.Version == params.GVR.Version) &&
			(r.Resource == "*" || r.Resource == params.GVR.Resource) {
			for _, verb := range r.Verbs {
				if verb == params.K8sVerb {
					ctrl.Log.V(0).Info("资源豁免组检查通过", "gvr", params.GVR)
					return true
				}
			}
		}
	}
	ctrl.Log.V(1).Info("资源豁免组检查未通过", "gvr", params.GVR)
	return false
}
