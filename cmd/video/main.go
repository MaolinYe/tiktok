package main

import (
	"log"
	video "tiktok/kitex/kitex_gen/video/videoservice"
)

func main() {
	svr := video.NewServer(new(VideoServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
