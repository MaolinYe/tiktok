package main

import (
	"net"
	"tiktok/kitex/kitex_gen/user/userservice"
	"tiktok/pkg/log"
	"tiktok/pkg/viper"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	logger = log.InitLogger("user")
)

func main() {
	config := viper.InitConfig()
	etcd_address := config.GetString("etcd.address")
	service_address := config.GetString("service.user.address")
	// 服务注册
	r, err := etcd.NewEtcdRegistry([]string{etcd_address})
	if err != nil {
		logger.Fatal(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", service_address)
	svr := userservice.NewServer(new(UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "user"}))
	// 运行服务
	err = svr.Run()
	if err != nil {
		logger.Println(err.Error())
	}
}
