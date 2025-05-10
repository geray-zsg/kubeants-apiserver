package main

import (
	"kubeants.io/config"
	"kubeants.io/initiallize"
)

func main() {
	// 先加载配置
	initiallize.Viper()

	// 根据配置初始化日志等级/格式
	initiallize.InitLogger()

	// 初始化其他组件
	initiallize.K8S()

	// 启动服务
	r := initiallize.Routers()
	r.Run(config.CONF.System.Port)
}
