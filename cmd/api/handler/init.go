package handler

import (
	"tiktok/pkg/log"
	"tiktok/pkg/viper"
)

var (
	Config = viper.InitConfig()
	logger = log.InitLogger("api")
)

func init() {
	InitUser()
	InitVideo()
	InitFavor()
	InitComment()
}
