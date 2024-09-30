package main

import (
	"log"
	favorite "tiktok/kitex/kitex_gen/favorite/favoriteservice"
)

func main() {
	svr := favorite.NewServer(new(FavoriteServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
