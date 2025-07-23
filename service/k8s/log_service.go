package k8s

import (
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	"kubeants.io/config"
)

type LogService struct{}

// GetPodLogs 获取指定 Pod 的日志
func (s *LogService) GetPodLogs(ctx context.Context, namespace, podName, container string, tailLines int64, follow bool) (string, error) {
	var err error

	// 日志选项
	opts := &corev1.PodLogOptions{
		Follow:     follow,
		TailLines:  &tailLines,
		Container:  container, // 如果为空，则取默认容器
		Timestamps: false,
	}

	req := config.KubeClientSet.CoreV1().Pods(namespace).GetLogs(podName, opts)

	readCloser, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("日志流获取失败: %v", err)
	}
	defer readCloser.Close()

	logs, err := io.ReadAll(readCloser)
	if err != nil {
		return "", fmt.Errorf("读取日志失败: %v", err)
	}

	return string(logs), nil
}
