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
	Cors   Cors   `yaml:"cors" json:"cors"`
}

/*
accessControlAllowCredentials: "true" # 是否允许浏览器发送跨域请求时携带认证信息（如Cookies和HTTP认证）
accessControlAllowHeaders: "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id" # 指定服务器接受哪些头部字段作为跨域请求的一部分
accessControlAllowMethods: "POST, GET, OPTIONS,DELETE,PUT" # 指定服务器允许的HTTP请求方法
accessControlExposeHeaders: "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At" # 指定哪些HTTP头部可以作为响应的一部分暴露给外部
*/
type Cors struct {
	Enable                        bool     `yaml:"enable" json:"enable"`
	DefaultOrigins                string   `yaml:"defaultOrigins" json:"defaultOrigins"`
	AllowedOrigins                []string `yaml:"allowedOrigins" json:"allowedOrigins"`
	AccessControlAllowCredentials string   `yaml:"accessControlAllowCredentials" json:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     string   `yaml:"accessControlAllowHeaders" json:"accessControlAllowHeaders"`
	AccessControlAllowMethods     string   `yaml:"accessControlAllowMethods" json:"accessControlAllowMethods"`
	AccessControlExposeHeaders    string   `yaml:"accessControlExposeHeaders" json:"accessControlExposeHeaders"`
}

var (
	CONF Server
	// 定义clientset全局变量用于任何地方都可以直接调用，再initiallize/k8s.go 中赋值
	KubeClientSet     *kubernetes.Clientset
	KubeDynamicClient dynamic.Interface
	Kubeconfig        string
)
