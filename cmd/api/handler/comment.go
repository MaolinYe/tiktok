package handler

import (
	"context"
	"net/http"
	"strconv"
	"tiktok/internal/response"
	"tiktok/kitex/kitex_gen/comment"
	"tiktok/kitex/kitex_gen/comment/commentservice"
	"tiktok/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	commentClient commentservice.Client
)

func InitComment() {
	etcd_address := Config.GetString("etcd.address")
	// 服务发现
	resolver, err := etcd.NewEtcdResolver([]string{etcd_address})
	if err != nil {
		logger.Fatal(err)
	}
	c, err := commentservice.NewClient("comment", client.WithResolver(resolver))
	if err != nil {
		logger.Fatal(err)
	}
	commentClient = c
	logger.Println("init comment client")
}

// 点赞操作
func CommentAction(ctx context.Context, c *app.RequestContext) {
	// 检查请求参数
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusOK, response.Base{
			StatusCode: -1,
			StatusMsg:  "action_type 不合法",
		})
		return
	}
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.Base{
			StatusCode: -1,
			StatusMsg:  "video_id 不合法",
		})
		return
	}
	token := c.Query("token")
	claims, _ := jwt.ValidateJWT(token)
	req := &comment.CommentActionRequest{
		VideoId:     videoID,
		ActionType:  int32(actionType),
		UserName:    claims.Username,
		CommentText: c.Query("comment_text"),
	}
	if actionType == 1 {
		logger.Println("comment video", videoID)
		if req.CommentText == "" {
			c.JSON(http.StatusOK, response.Base{
				StatusCode: -1,
				StatusMsg:  "评论为空",
			})
			return
		}
	} else {
		logger.Println("cancel comment video", videoID)
		commentID, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, response.Base{
				StatusCode: -1,
				StatusMsg:  "comment_id 不合法",
			})
			return
		}
		req.CommentId = commentID
	}
	res, _ := commentClient.CommentAction(ctx, req)
	// 检查服务是否上线
	if res == nil {
		logger.Println("comment无服务")
		c.JSON(http.StatusOK, response.Base{
			StatusCode: -1,
			StatusMsg:  "无服务",
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Base{
			StatusCode: -1,
			StatusMsg:  res.StatusMsg,
		})
		return
	}
	c.JSON(http.StatusOK, response.CommentAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		}, Comment: res.Comment,
	})
}

// 获取点赞视频列表
func CommentList(ctx context.Context, c *app.RequestContext) {
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.Base{
			StatusCode: -1,
			StatusMsg:  "video_id 不合法",
		})
		return
	}
	logger.Println("comment list video", videoID)
	req := &comment.CommentListRequest{VideoId: videoID}
	res, _ := commentClient.CommentList(ctx, req)
	// 检查服务是否上线
	if res == nil {
		logger.Println("comment无服务")
		c.JSON(http.StatusOK, response.Base{
			StatusCode: -1,
			StatusMsg:  "无服务",
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Base{
			StatusCode: -1,
			StatusMsg:  res.StatusMsg,
		})
		return
	}
	c.JSON(http.StatusOK, response.CommentList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		CommentList: res.CommentList,
	})
}
