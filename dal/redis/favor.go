package redis

import (
	"context"
	"fmt"
	"strconv"
	"tiktok/dal/db"

	"github.com/go-redis/redis/v8"
)

// 点赞视频
func FavorVideo(ctx context.Context, videoID int64) error {
	rdb := rdb_favor
	key := fmt.Sprintf("%d", videoID)
	// 重置过期时间
	exits, err := rdb.Expire(ctx, key, expireTime).Result()
	if err != nil {
		logger.Println(err)
		return err
	}
	// 不存在读取MySQL再自增
	if !exits {
		// 上互斥锁
		mutex_favor.Lock()
		defer mutex_favor.Unlock()
		// 判断是否有同行新增键了
		had, err := rdb.Exists(ctx, key).Result()
		if err != nil {
			logger.Println(err)
			return err
		}
		if had == 0 {
			// 读取MySQL
			favorCount, err := db.GetFavorCount(ctx, videoID)
			if err != nil {
				logger.Println(err)
				return err
			}
			// 新增key自增
			favorCount++
			if err = rdb.Set(ctx, key, favorCount, expireTime).Err(); err != nil {
				logger.Println(err)
				return err
			}
			return nil
		}
	}
	// 存在直接自增
	if err := rdb.Incr(ctx, key).Err(); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

// 取消点赞
func CancelFavorVideo(ctx context.Context, videoID int64) error {
	rdb := rdb_favor
	key := fmt.Sprintf("%d", videoID)
	// 重置过期时间
	exits, err := rdb.Expire(ctx, key, expireTime).Result()
	if err != nil {
		logger.Println(err)
		return err
	}
	// 不存在读取MySQL再自减
	if !exits {
		// 上互斥锁
		mutex_favor.Lock()
		defer mutex_favor.Unlock()
		// 判断是否有同行新增键了
		had, err := rdb.Exists(ctx, key).Result()
		if err != nil {
			logger.Println(err)
			return err
		}
		if had == 0 {
			// 读取MySQL
			favorCount, err := db.GetFavorCount(ctx, videoID)
			if err != nil {
				logger.Println(err)
				return err
			}
			// 新增key自减
			favorCount--
			if err = rdb.Set(ctx, key, favorCount, expireTime).Err(); err != nil {
				logger.Println(err)
				return err
			}
			return nil
		}
	}
	// 存在直接自减
	if err := rdb.Decr(ctx, key).Err(); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

// 同步数据到MySQL
func SyncFavorToMySQL(ctx context.Context) {
	rdb := rdb_favor
	var cursor uint64 = 0
	for {
		// 使用 SCAN 命令获取键
		keys, newCursor, err := rdb.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			logger.Println("Error scanning keys:", err)
			return
		}
		// 处理当前批次的键
		for _, key := range keys {
			value, err := rdb.Get(ctx, key).Result()
			if err != nil && err != redis.Nil {
				logger.Println(err)
				continue
			}
			logger.Printf("Key: %s, Value: %s\n", key, value)
			favorCount, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				logger.Println(err)
				continue
			}
			videoID, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				logger.Println(err)
				continue
			}
			err=db.UpdateFavorCount(ctx, videoID, favorCount)
			if err != nil {
				logger.Println("同步失败", err)
				return
			}
		}
		// 更新光标
		cursor = newCursor
		if cursor == 0 {
			break // 直到扫描完所有的键
		}
	}
}
