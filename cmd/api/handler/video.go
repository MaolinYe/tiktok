package handler

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"tiktok/internal/response"
	"tiktok/kitex/kitex_gen/video"
	"tiktok/kitex/kitex_gen/video/videoservice"
	"tiktok/pkg/jwt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	videoClient videoservice.Client
)

func InitVideo() {
	etcd_address := Config.GetString("etcd.address")
	// 服务发现
	resolver, err := etcd.NewEtcdResolver([]string{etcd_address})
	if err != nil {
		logger.Fatal(err)
	}
	c, err := videoservice.NewClient("video", client.WithResolver(resolver))
	if err != nil {
		logger.Fatal(err)
	}
	videoClient = c
	logger.Println("init video client")
}

// 获取视频流
func Feed(ctx context.Context, c *app.RequestContext) {
	latestTime := c.Query("latest_time")
	logger.Println("feed", latestTime)
	var timestamp int64 = 0
	if latestTime != "" {
		timestamp, _ = strconv.ParseInt(latestTime, 10, 64)
	} else {
		timestamp = time.Now().UnixMilli()
	}
	req := &video.FeedRequest{
		LatestTime: timestamp,
	}
	res, _ := videoClient.Feed(ctx, req)
	if res == nil {
		logger.Println("video无服务")
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "无服务",
			},
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Feed{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.Feed{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: res.VideoList,
	})
}

// 获取视频发布列表
func PublishList(ctx context.Context, c *app.RequestContext) {
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	logger.Println("publishList", uid)
	if err != nil {
		c.JSON(http.StatusOK, response.PublishList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
		})
		return
	}
	req := &video.PublishListRequest{
		UserId: uid,
	}
	res, _ := videoClient.PublishList(ctx, req)
	if res == nil {
		logger.Println("video无服务")
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "无服务",
			},
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.PublishList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: res.VideoList,
	})
}

// 投稿视频
func PublishAction(ctx context.Context, c *app.RequestContext) {
	title := c.PostForm("title")
	token := c.PostForm("token")
	claims, _ := jwt.ValidateJWT(token)
	user_name := claims.Username
	log.Println("publish action", user_name)
	if title == "" {
		c.JSON(http.StatusBadRequest, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "标题不能为空",
			},
		})
		return
	}
	// 视频数据
	file, err := c.FormFile("data")
	if err != nil {
		logger.Println(err.Error())
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			},
		})
		return
	}
	// 视频文件转二进制
	src, err := file.Open()
	if err != nil {
		logger.Println(err.Error())
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			},
		})
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		logger.Println(err.Error())
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			},
		})
		return
	}
	req := &video.PublishActionRequest{
		Title:    title,
		Data:     buf.Bytes(),
		UserName: user_name,
	}
	res, _ := videoClient.PublishAction(ctx, req)
	if res == nil {
		logger.Println("video无服务")
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "无服务",
			},
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.PublishAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success",
		},
	})
}
