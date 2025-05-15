package config

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type System struct {
	Port string `json:"port" yaml:"port" mapstructure:"port"`
}

type JWT struct {
	Secret     string `json:"secret" yaml:"secret" mapstructure:"secret"`
	Expiration int    `json:"expiration" yaml:"expiration" mapstructure:"expiration"`
}

type Log struct {
	Level  string `mapstructure:"level" json:"level" yaml:"level"`    // debug, info, warn, error
	Format string `mapstructure:"format" json:"format" yaml:"format"` // console 开发 或 json
	File   string `mapstructure:"file" json:"file" yaml:"file"`       // 日志写入文件路径
}

// 权限控制
type ExemptResource struct {
	Group    string   `mapstructure:"group" json:"group" yaml:"group"`
	Version  string   `mapstructure:"version" json:"version" yaml:"version"`
	Resource string   `mapstructure:"resource" json:"resource" yaml:"resource"`
	Verbs    []string `mapstructure:"verbs" json:"verbs" yaml:"verbs"`
}

type Authz struct {
	ExemptResources []ExemptResource `mapstructure:"exemptResources" json:"exemptResources" yaml:"exemptResources"`
}

type Server struct {
	System System `mapstructure:"system"`
	JWT    JWT    `mapstructure:"jwt"`
	Cors   Cors   `mapstructure:"cors"`
	Log    Log    `mapstructure:"log"`
	Authz  Authz  `mapstructure:"authz"`
}

/*
accessControlAllowCredentials: "true" # 是否允许浏览器发送跨域请求时携带认证信息（如Cookies和HTTP认证）
accessControlAllowHeaders: "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id" # 指定服务器接受哪些头部字段作为跨域请求的一部分
accessControlAllowMethods: "POST, GET, OPTIONS,DELETE,PUT" # 指定服务器允许的HTTP请求方法
accessControlExposeHeaders: "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At" # 指定哪些HTTP头部可以作为响应的一部分暴露给外部
*/
type Cors struct {
	Enable                        bool     `mapstructure:"enable"`
	DefaultOrigins                string   `mapstructure:"defaultOrigins"`
	AllowedOrigins                []string `mapstructure:"allowedOrigins"`
	AccessControlAllowCredentials string   `mapstructure:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     string   `mapstructure:"accessControlAllowHeaders"`
	AccessControlAllowMethods     string   `mapstructure:"accessControlAllowMethods"`
	AccessControlExposeHeaders    string   `mapstructure:"accessControlExposeHeaders"`
}

var (
	CONF Server
	// 定义clientset全局变量用于任何地方都可以直接调用，再initiallize/k8s.go 中赋值
	KubeClientSet     *kubernetes.Clientset
	KubeDynamicClient dynamic.Interface
	Kubeconfig        string
)
