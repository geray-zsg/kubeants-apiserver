package initiallize

import (
	"fmt"

	"github.com/spf13/viper"
	"kubeants.io/config"
)

func Viper() {
	v := viper.New()

	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	// 结构体字段需使用 mapstructure 标签，防止大小写解析错误
	if err := v.Unmarshal(&config.CONF); err != nil {
		panic(fmt.Errorf("配置绑定失败: %w", err))
	}
}
