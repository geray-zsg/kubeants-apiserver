package pod

import (
	"context"
	"fmt"
	"time"

	"kubeant.cn/config"

	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodService struct{}

// Get a pod detail in namespace
func (*PodService) GetoldPodDetail(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	getpod, err := config.KubeClientSet.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return getpod, nil
}

// update pod (delete+create)
func (*PodService) UpdatePod(ctx context.Context, k8sPod *corev1.Pod) (*corev1.Pod, error) {

	updatePod, err := config.KubeClientSet.CoreV1().Pods(k8sPod.Namespace).Update(ctx, k8sPod, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return updatePod, nil
}

func (p *PodService) DeleteAndCreatePod(ctx context.Context, k8sPod *corev1.Pod) (*corev1.Pod, error) {
	clientSetPod := config.KubeClientSet.CoreV1().Pods(k8sPod.Namespace)
	geteOldPod, err := clientSetPod.Get(ctx, k8sPod.Name, metav1.GetOptions{})
	if err == nil {
		// Pod存在 删除后重建
		// delete old Pod
		// err := clientSetPod.Delete(ctx, k8sPod.Name, metav1.DeleteOptions{})
		err := p.DeletePod(ctx, k8sPod.Namespace, k8sPod.Name)
		if err != nil {
			return nil, err
		}

		// 设置一个超时时间，防止无限等待
		deadline := time.Now().Add(time.Second * 30)
		// 轮询直到Pod被删除
		for {
			_, err := clientSetPod.Get(ctx, k8sPod.Name, metav1.GetOptions{})
			if k8serror.IsNotFound(err) {
				break // Pod已经被删除
			}
			if err != nil {
				return nil, err // 其他错误
			}
			if time.Now().After(deadline) {
				return nil, fmt.Errorf("failed to update Pod,because time out waiting for old pod[%s] to be deleted", k8sPod.Name)
			}

			time.Sleep(time.Second * 1)
		}

	} else if !k8serror.IsNotFound(err) {
		// 其他错误
		fmt.Println(geteOldPod)
		return nil, err
	}

	createPod, err := clientSetPod.Create(ctx, k8sPod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return createPod, nil
}

func (*PodService) DeletePod(ctx context.Context, namespace, name string) error {
	propagationPolicy := metav1.DeletePropagationBackground // 后台删除
	gracePeriodSeconds := int64(0)
	return config.KubeClientSet.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
		PropagationPolicy:  &propagationPolicy,
	})
}
