package namespace

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeants.com/config"
)

type NamespaceService struct{}

func (*NamespaceService) GetNamespaceDetail(ctx context.Context, name string) (*corev1.Namespace, error) {
	namespaceDetail, err := config.KubeClientSet.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return namespaceDetail, nil
}

func (*NamespaceService) DeleteNamespace(ctx context.Context, name string) error {
	err := config.KubeClientSet.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (*NamespaceService) CreateNamespace(ctx context.Context, k8sNamespace *corev1.Namespace) (*corev1.Namespace, error) {

	createNamespace, err := config.KubeClientSet.CoreV1().Namespaces().Create(ctx, k8sNamespace, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return createNamespace, nil
}
