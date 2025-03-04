package initiallize

import (
	"github.com/spf13/viper"
	"kubeants.io/config"
)

func Viper() {
	v := viper.New()

	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}

	// 参数实体绑定
	v.Unmarshal(&config.CONF)
}
