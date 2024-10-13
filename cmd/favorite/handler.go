package main

import (
	"context"
	"errors"
	"tiktok/dal/db"
	"tiktok/dal/redis"
	"tiktok/kitex/kitex_gen/favorite"
	"tiktok/kitex/kitex_gen/user"
	"tiktok/kitex/kitex_gen/video"
	"tiktok/pkg/minio"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	logger.Println("favoring...")
	// 获取用户id
	user, err := db.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		resp = &favorite.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	} else if user == nil {
		resp = &favorite.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}
		logger.Println(resp.StatusMsg)
		return resp, nil
	}
	favor := &db.Favor{VideoID: uint(req.VideoId), UserID: uint(user.ID)}
	if req.ActionType == 1 { // 点赞
		// 创建点赞记录
		if err := db.CreateFavor(ctx, favor); err != nil {
			resp = &favorite.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		}
		// 更新点赞数
		if err := redis.FavorVideo(ctx, int64(favor.VideoID)); err != nil {
			resp = &favorite.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		}
	} else { // 取消点赞
		// 删除点赞记录
		if err := db.DeleteFavor(ctx, favor); err != nil {
			resp = &favorite.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		}
		// 更新点赞数
		if err := redis.CancelFavorVideo(ctx, int64(favor.VideoID)); err != nil {
			resp = &favorite.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误",
			}
			logger.Println(resp.StatusMsg, err)
			return resp, nil
		}
	}
	resp = &favorite.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return resp, nil
}

// FavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	logger.Println("favorite list getting...")
	// 获取用户视频发布列表
	videos, err := db.GetFavorList(ctx, req.UserId)
	if err != nil {
		resp = &favorite.FavoriteListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		logger.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 构造视频信息返回
	list, err := getVideoToResponse(ctx, videos, req.UserId)
	if err != nil {
		resp = &favorite.FavoriteListResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}
		logger.Println(err)
		return resp, nil
	}
	resp = &favorite.FavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  list,
	}
	return resp, nil
}

// from db to response
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
