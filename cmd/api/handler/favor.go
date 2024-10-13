package handler

import (
	"context"
	"net/http"
	"strconv"
	"tiktok/internal/response"
	"tiktok/kitex/kitex_gen/favorite"
	"tiktok/kitex/kitex_gen/favorite/favoriteservice"
	"tiktok/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	favorClient favoriteservice.Client
)

func InitFavor() {
	etcd_address := Config.GetString("etcd.address")
	// 服务发现
	resolver, err := etcd.NewEtcdResolver([]string{etcd_address})
	if err != nil {
		logger.Fatal(err)
	}
	c, err := favoriteservice.NewClient("favor", client.WithResolver(resolver))
	if err != nil {
		logger.Fatal(err)
	}
	favorClient = c
	logger.Println("init favor client")
}

// 点赞操作
func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusOK, response.FavoriteAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "action_type 不合法",
			},
		})
		return
	}
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.FavoriteAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
		})
		return
	}
	if actionType == 1 {
		logger.Println("favor video", videoID)
	} else {
		logger.Println("cancel favor video", videoID)
	}
	token := c.Query("token")
	claims, _ := jwt.ValidateJWT(token)
	req := &favorite.FavoriteActionRequest{
		UserName: claims.Username,
		VideoId:    videoID,
		ActionType: int32(actionType),
	}
	res, _ := favorClient.FavoriteAction(ctx, req)
	// 检查服务是否上线
	if res == nil {
		logger.Println("favor无服务")
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "无服务",
			},
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.FavoriteAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.FavoriteAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
	})
}

// 获取点赞视频列表
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	logger.Println("favorite list user", userID)
	if err != nil {
		c.JSON(http.StatusOK, response.FavoriteList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
		})
		return
	}
	req := &favorite.FavoriteListRequest{
		UserId: userID,
	}
	res, _ := favorClient.FavoriteList(ctx, req)
	// 检查服务是否上线
	if res == nil {
		logger.Println("favor无服务")
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "无服务",
			},
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.FavoriteList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.FavoriteList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		VideoList: res.VideoList,
	})
}
