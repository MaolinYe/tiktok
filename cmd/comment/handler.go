package main

import (
	"context"
	"encoding/json"
	"errors"
	"tiktok/dal/db"
	"tiktok/dal/redis"
	"tiktok/kitex/kitex_gen/comment"
	"tiktok/kitex/kitex_gen/user"
	"tiktok/pkg/minio"
	"time"

	"github.com/streadway/amqp"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	// 获取用户id
	user, err := db.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		resp = &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	} else if user == nil {
		resp = &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}
		logger.Println(resp.StatusMsg)
		return resp, nil
	}
	// 发布评论
	if req.ActionType == 1 {
		logger.Println("commenting...", req.VideoId)
		remark := &db.Comment{VideoID: uint(req.VideoId), UserID: user.ID, Content: req.CommentText}
		// 创建评论记录
		msg, err := json.Marshal(remark)
		if err != nil {
			resp = &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		}
		mq.PublishSimple(msg)
		// 更新视频评论数
		if err = redis.CommentVideo(ctx, req.VideoId); err != nil {
			resp = &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		}
		// 获取作者信息
		author, err := getAuthorData(ctx, int64(user.ID))
		if err != nil {
			resp = &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		}
		resp = &comment.CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			Comment: &comment.Comment{
				Id:         int64(remark.ID),
				User:       author,
				Content:    remark.Content,
				CreateDate: time.Now().Format("2006-01-02 15:04:05"),
			},
		}
		return resp, nil
	}
	// 删除评论
	logger.Println("comment deleting...", req.VideoId)
	// 获取视频作者
	video, err := db.GetVideoByID(ctx, req.VideoId)
	if err != nil {
		resp = &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 获取评论作者
	remark, err := db.GetCommentByID(ctx, req.CommentId)
	if err != nil {
		resp = &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 视频作者或者评论作者可以删除
	if user.ID != video.Author.ID && user.ID != remark.UserID {
		resp = &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "无权限",
		}
		logger.Println(resp.StatusMsg)
		return resp, nil
	}
	// 删除评论记录
	if err = db.DeleteComment(ctx, req.CommentId); err != nil {
		resp = &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 更新视频评论数
	if err = redis.DeleteCommentVideo(ctx, req.VideoId); err != nil {
		resp = &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	resp = &comment.CommentActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	logger.Println("comment list getting...")
	comments, err := db.GetCommentList(ctx, req.VideoId)
	if err != nil {
		resp = &comment.CommentListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	list := make([]*comment.Comment, len(comments))
	for i := 0; i < len(comments); i++ {
		author, err := getAuthorData(ctx, int64(comments[i].UserID))
		if err != nil {
			resp = &comment.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		}
		list[i] = &comment.Comment{
			Id:         int64(comments[i].ID),
			User:       author,
			Content:    comments[i].Content,
			CreateDate: comments[i].CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	resp = &comment.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		CommentList: list,
	}
	return
}

// 获取作者信息
func getAuthorData(ctx context.Context, authorID int64) (data *user.User, err error) {
	// 获取作者
	author, err := db.GetUserByID(ctx, authorID)
	if err != nil {
		return nil, errors.New("服务器错误")
	} else if author == nil {
		return nil, errors.New("用户不存在")
	}
	// 获取头像
	avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, author.Avatar)
	if err != nil {
		return nil, errors.New("服务器内部错误：获取头像失败")
	}
	// 获取背景
	backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, author.BackgroundImage)
	if err != nil {
		return nil, errors.New("服务器内部错误：获取背景图失败")
	}
	data = &user.User{
		Id:              int64(author.ID),
		Name:            author.UserName,
		FollowerCount:   int64(author.FollowerCount),
		FollowCount:     int64(author.FollowingCount),
		IsFollow:        true,
		Avatar:          avatar,
		BackgroundImage: backgroundImage,
		Signature:       author.Signature,
		TotalFavorited:  int64(author.TotalFavorited),
		WorkCount:       int64(author.WorkCount),
		FavoriteCount:   int64(author.FavoriteCount),
	}
	return
}

func consume(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		comment := new(db.Comment)
		if err := json.Unmarshal(msg.Body, comment); err != nil {
			logger.Println(err)
			continue
		}
		if err := db.CreateComment(context.Background(), comment); err != nil {
			logger.Println(err)
			continue
		}
	}
}
