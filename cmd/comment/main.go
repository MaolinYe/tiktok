package main

import (
	"context"
	"fmt"
	"net"
	"tiktok/dal/redis"
	"tiktok/kitex/kitex_gen/comment/commentservice"
	"tiktok/pkg/log"
	"tiktok/pkg/viper"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/robfig/cron"
)

var (
	logger = log.InitLogger("comment")
	config = viper.InitConfig()
)

func main() {
	etcd_address := config.GetString("etcd.address")
	service_address := config.GetString("service.comment.address")
	// 服务注册
	r, err := etcd.NewEtcdRegistry([]string{etcd_address})
	if err != nil {
		logger.Fatal(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", service_address)
	svr := commentservice.NewServer(new(CommentServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "comment"}))
	// 启动定时任务同步MySQL
	frequency := fmt.Sprintf("@every %dm", redis.SyncTime)
	c := cron.New()
	c.AddFunc(frequency, func() {
		redis.SyncCommentToMySQL(context.Background())
	})
	c.Start()
	// 运行服务
	err = svr.Run()
	if err != nil {
		logger.Println(err.Error())
	}
}
