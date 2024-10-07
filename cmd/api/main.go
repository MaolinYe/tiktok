package main

import (
	"log"

	"tiktok/cmd/api/handler"
	"tiktok/pkg/jwt"
	"tiktok/pkg/logger"
	"tiktok/pkg/viper"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func init() {
	logger.InitLogger()
}

func main() {
	config := viper.InitConfig()
	api_address := config.GetString("api.address")
	hz := server.New(server.WithHostPorts(api_address))
	registerGroup(hz)
	if err := hz.Run(); err != nil {
		log.Fatal(err)
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
		// publish := douyin.Group("/publish")
		// {
		// 	publish.GET("/list/", handler.PublishList)
		// 	publish.POST("/action/", handler.PublishAction)
		// }
		// douyin.GET("/feed", handler.Feed)
		// favorite := douyin.Group("/favorite")
		// {
		// 	favorite.POST("/action/", handler.FavoriteAction)
		// 	favorite.GET("/list/", handler.FavoriteList)
		// }
		// comment := douyin.Group("/comment")
		// {
		// 	comment.POST("/action/", handler.CommentAction)
		// 	comment.GET("/list/", handler.CommentList)
		// }
	}
}
