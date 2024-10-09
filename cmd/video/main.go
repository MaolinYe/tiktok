package main

import (
	"net"
	"tiktok/kitex/kitex_gen/video/videoservice"
	"tiktok/pkg/log"
	"tiktok/pkg/viper"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	logger = log.InitLogger("video")
	config = viper.InitConfig()
	limit  = config.GetInt("service.video.limit")
)

func main() {
	etcd_address := config.GetString("etcd.address")
	service_address := config.GetString("service.video.address")
	// 服务注册
	r, err := etcd.NewEtcdRegistry([]string{etcd_address})
	if err != nil {
		logger.Fatal(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", service_address)
	svr := videoservice.NewServer(new(VideoServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "video"}))
	// 运行服务
	err = svr.Run()
	if err != nil {
		logger.Println(err.Error())
	}
}
