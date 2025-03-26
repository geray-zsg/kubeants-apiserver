package main

import (
	"kubeants.io/config"
	"kubeants.io/initiallize"
	logger "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	// 初始化日志
	logger.SetLogger(zap.New(zap.UseDevMode(true))) // 初始化 Zap 日志
	r := initiallize.Routers()
	// 初始化参数
	initiallize.Viper()
	initiallize.K8S()

	r.Run(config.CONF.System.Port)
}
