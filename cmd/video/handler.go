package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"tiktok/dal/db"
	"tiktok/kitex/kitex_gen/user"
	"tiktok/kitex/kitex_gen/video"
	"tiktok/pkg/ffmpeg"
	"tiktok/pkg/minio"
	"time"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	logger.Println("feeding...")
	// 获取用户视频发布列表
	videos, err := db.GetFeed(ctx, limit, req.LatestTime)
	if err != nil {
		resp = &video.FeedResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	var user_id int64
	// 如果在登录状态获取用户id
	if req.UserName != "" {
		user, err := db.GetUserByUsername(ctx, req.UserName)
		if err != nil {
			resp = &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		} else if user == nil {
			resp = &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			}
			logger.Println(resp.StatusMsg)
			return resp, nil
		}
		user_id = int64(user.ID)
	}
	// 构造视频信息返回
	list, err := getVideoToResponse(ctx, videos, user_id)
	if err != nil {
		resp = &video.FeedResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}
		logger.Println(err)
		return resp, nil
	}
	nextTime := time.Now().UnixMilli()
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].CreatedAt.UnixMilli()
	}
	resp = &video.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  list,
		NextTime:   nextTime,
	}
	return resp, nil
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	logger.Println("publishing...")
	// 获取用户id
	user, err := db.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		resp = &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	} else if user == nil {
		resp = &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}
		logger.Println(resp.StatusMsg)
		return resp, nil
	}
	// 上传视频
	videoName := fmt.Sprintf("%d_%d.mp4", user.ID, time.Now().UnixMilli())
	videoSize := len(req.Data)
	logger.Println("video size:", videoSize)
	uploadSize, err := minio.UploadFileByIO(minio.VideoBucketName, videoName, bytes.NewReader(req.Data), int64(videoSize), "application/mp4")
	if err != nil {
		resp = &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：上传视频失败",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	logger.Println("upload size:", uploadSize)
	videoURL, err := minio.GetFileTemporaryURL(minio.VideoBucketName, videoName)
	if err != nil {
		resp = &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：获取视频URL失败",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 获取封面
	coverData, err := ffmpeg.TakeFrameFromVideo(videoURL, 1)
	if err != nil {
		resp = &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：获取视频封面失败",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 上传封面
	coverName := fmt.Sprintf("%d_%d.jpg", user.ID, time.Now().UnixMilli())
	coverSize := coverData.Len()
	logger.Println("cover size:", coverSize)
	uploadSize, err = minio.UploadFileByIO(minio.CoverBucketName, coverName, coverData, int64(coverSize), "image/png")
	if err != nil {
		resp = &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：上传视频封面失败",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	logger.Println("upload size:", uploadSize)
	// 写入数据库
	v := &db.Video{
		Title:    req.Title,
		PlayUrl:  videoName,
		CoverUrl: coverName,
		AuthorID: user.ID,
	}
	if err = db.CreateVideo(ctx, v); err != nil {
		resp = &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	resp = &video.PublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return resp, nil
}

// PublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	logger.Println("publishing list getting...")
	// 获取用户视频发布列表
	videos, err := db.GetPublishList(ctx, req.UserId)
	if err != nil {
		resp = &video.PublishListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 构造视频信息返回
	list, err := getVideoToResponse(ctx, videos, req.UserId)
	if err != nil {
		resp = &video.PublishListResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}
		logger.Println(err)
		return resp, nil
	}
	resp = &video.PublishListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  list,
	}
	return resp, nil
}

func getVideoToResponse(ctx context.Context, videos []*db.Video, userID int64) (list []*video.Video, err error) {
	for i := 0; i < len(videos); i++ {
		// 获取作者
		author, err := db.GetUserByID(ctx, int64(videos[i].AuthorID))
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
		// 获取视频封面
		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, videos[i].PlayUrl)
		if err != nil {
			return nil, errors.New("服务器内部错误：获取头像失败")
		}
		// 获取视频
		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, videos[i].CoverUrl)
		if err != nil {
			return nil, errors.New("服务器内部错误：获取视频失败")
		}
		// 是否已点赞
		isFavorite, err := db.IsFavorite(ctx, videos[i].ID, userID)
		if err != nil {
			return nil, errors.New("服务器内部错误")
		}
		list = append(list, &video.Video{
			Id: int64(videos[i].ID),
			Author: &user.User{
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
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(videos[i].FavoriteCount),
			CommentCount:  int64(videos[i].CommentCount),
			IsFavorite:    isFavorite,
			Title:         videos[i].Title,
		})
	}
	return list, nil
}
