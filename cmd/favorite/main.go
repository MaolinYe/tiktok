package main

import (
	"context"
	"fmt"
	"net"
	"tiktok/dal/redis"
	"tiktok/kitex/kitex_gen/favorite/favoriteservice"
	"tiktok/pkg/log"
	"tiktok/pkg/viper"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/robfig/cron"
)

var (
	logger = log.InitLogger("favor")
	config = viper.InitConfig()
)

func main() {
	etcd_address := config.GetString("etcd.address")
	service_address := config.GetString("service.favor.address")
	// 服务注册
	r, err := etcd.NewEtcdRegistry([]string{etcd_address})
	if err != nil {
		logger.Fatal(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", service_address)
	svr := favoriteservice.NewServer(new(FavoriteServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "favor"}))
	// 启动定时任务同步MySQL
	frequency := fmt.Sprintf("@every %dm", redis.SyncTime)
	c := cron.New()
	c.AddFunc(frequency, func() {
		redis.SyncFavorToMySQL(context.Background())
	})
	c.Start()
	// 运行服务
	err = svr.Run()
	if err != nil {
		logger.Println(err.Error())
	}
}
