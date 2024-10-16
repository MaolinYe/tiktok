package redis

import (
	"sync"
	"tiktok/pkg/log"
	"tiktok/pkg/viper"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb_favor   *redis.Client
	rdb_comment *redis.Client
	config      = viper.InitConfig()
	logger      = log.InitLogger("redis")
	expireTime  = time.Duration(config.GetUint("redis.expireTime")) * time.Minute
	SyncTime    = config.GetUint("redis.syncTime")
	mutex_favor sync.Mutex
	mutex_comment sync.Mutex
)

func init() {
	// 创建 Redis 客户端 点赞
	rdb_favor = redis.NewClient(&redis.Options{
		Addr:     config.GetString("redis.address"),  // Redis 服务器地址
		Password: config.GetString("redis.password"), // Redis 密码，如果没有则留空
		DB:       int(config.GetUint("redis.favorDB")),
	})
	if rdb_favor == nil {
		panic("fail to connect redis")
	}
	// 创建 Redis 客户端 评论
	rdb_comment = redis.NewClient(&redis.Options{
		Addr:     config.GetString("redis.address"),  // Redis 服务器地址
		Password: config.GetString("redis.password"), // Redis 密码，如果没有则留空
		DB:       int(config.GetUint("redis.commentDB")),
	})
	if rdb_comment == nil {
		panic("fail to connect redis")
	}
	logger.Println("redis connected")
}
