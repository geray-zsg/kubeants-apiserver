package main

import (
	"kubeant.cn/config"
	"kubeant.cn/initiallize"
)

func main() {
	r := initiallize.Routers()
	// 初始化参数
	initiallize.Viper()
	initiallize.K8S()

	r.Run(config.CONF.System.Port)
}
