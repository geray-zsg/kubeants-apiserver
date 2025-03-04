package main

import (
	"kubeants.io/config"
	"kubeants.io/initiallize"
)

func main() {
	r := initiallize.Routers()
	// 初始化参数
	initiallize.Viper()
	initiallize.K8S()

	r.Run(config.CONF.System.Port)
}
