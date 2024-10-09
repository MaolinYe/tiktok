package main

import (
	"tiktok/cmd/api/handler"
	"tiktok/pkg/jwt"
	"tiktok/pkg/viper"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	config := viper.InitConfig()
	api_address := config.GetString("api.address")
	maxRequestBodySize := config.GetInt("api.maxRequestBodySize")
	hz := server.New(server.WithHostPorts(api_address), server.WithMaxRequestBodySize(maxRequestBodySize))
	registerGroup(hz)
	if err := hz.Run(); err != nil {
		panic(err)
	}
}

func registerGroup(hz *server.Hertz) {
	douyin := hz.Group("/douyin")
	{
		user := douyin.Group("/user")
		{
			user.GET("/", jwt.AuthMiddleware, handler.UserInfo)
			user.POST("/register/", handler.Register)
			user.POST("/login/", handler.Login)
		}
		publish := douyin.Group("/publish")
		{
			publish.GET("/list/", jwt.AuthMiddleware, handler.PublishList)
			publish.POST("/action/", jwt.AuthMiddleware, handler.PublishAction)
		}
		douyin.GET("/feed", handler.Feed)
	}
}
