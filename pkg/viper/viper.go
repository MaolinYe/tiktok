package viper

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig() viper.Viper {
	config := viper.New()
	// 设置配置文件的名称和类型
	config.SetConfigName("config")       // 不带扩展名
	config.SetConfigType("yaml")         // 设置配置文件类型
	config.AddConfigPath("../../config/") // 添加配置文件所在路径

	// 读取配置文件
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	return *config
}
