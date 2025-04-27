package config

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type System struct {
	Port string `json:"port" yaml:"port"`
}

type JWT struct {
	Secret     string `json:"secret" yaml:"secret"`
	Expiration int    `json:"expiration" yaml:"expiration"`
}

type Server struct {
	System System `json:"System" yaml:"system"`
	JWT    JWT    `json:"JWT" yaml:"jwt"`
}

var (
	CONF Server
	// 定义clientset全局变量用于任何地方都可以直接调用，再initiallize/k8s.go 中赋值
	KubeClientSet     *kubernetes.Clientset
	KubeDynamicClient dynamic.Interface
	Kubeconfig        string
)
