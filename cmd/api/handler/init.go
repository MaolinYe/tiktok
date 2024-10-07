package handler

import (
	"tiktok/pkg/viper"
)

var (
	Config = viper.InitConfig()
)

func init() {
	InitUser()
}
