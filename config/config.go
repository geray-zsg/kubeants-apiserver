package config

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type System struct {
	Port string `json:"port" yaml:"port"`
}

type Server struct {
	System System `json:"System" yaml:"system"`
}

var (
	CONF Server
	// 定义clientset全局变量用于任何地方都可以直接调用，再initiallize/k8s.go 中赋值
	KubeClientSet     *kubernetes.Clientset
	KubeDynamicClient dynamic.Interface
)
