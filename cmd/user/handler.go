package main

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
	"tiktok/dal/db"
	"tiktok/kitex/kitex_gen/user"
	"tiktok/pkg/jwt"
	"tiktok/pkg/minio"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	log.Println("registering...")
	// 检查用户名是否冲突
	usr, err := db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Println(err)
		resp = &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		log.Println(resp.StatusMsg, err)
		return resp, nil
	} else if usr != nil {
		resp = &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "用户名已存在",
		}
		log.Println(resp.StatusMsg)
		return resp, nil
	}
	// 生成token
	token, err := jwt.GenerateJWT(req.Username)
	if err != nil {
		log.Println(err)
		resp = &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		log.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		resp = &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		log.Println(resp.StatusMsg, err)
		return resp, nil
	}
	// 创建用户
	usr = &db.User{
		UserName: req.Username,
		Password: string(hashedPassword),
	}
	if err = db.CreateUser(ctx, usr); err != nil {
		log.Println(err)
		resp = &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		log.Println(resp.StatusMsg, err)
		return resp, nil
	}
	resp = &user.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(usr.ID),
		Token:      token,
	}
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	log.Println("loginning...")
	// 检查用户是否存在
	usr, err := db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Println(err)
		resp = &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		log.Println(resp.StatusMsg, err)
		return resp, nil
	} else if usr == nil {
		resp = &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}
		log.Println(resp.StatusMsg)
		return resp, nil
	}
	// 匹配密码
	if bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password)) != nil {
		resp = &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "密码错误",
		}
		log.Println(resp.StatusMsg)
		return resp, nil
	}
	// 生成token
	token, err := jwt.GenerateJWT(req.Username)
	if err != nil {
		log.Println(err)
		resp = &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		log.Println(resp.StatusMsg, err)
		return resp, nil
	}
	resp = &user.UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(usr.ID),
		Token:      token,
	}
	return resp, nil
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	log.Println("userinfo getting...")
	// 检查用户是否存在
	usr, err := db.GetUserByID(ctx, req.UserId)
	if err != nil {
		log.Println(err)
		resp = &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误",
		}
		log.Println(resp.StatusMsg)
		return resp, nil
	} else if usr == nil {
		resp = &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}
		log.Println(resp.StatusMsg)
		return resp, nil
	}
	// 获取头像
	avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
	if err != nil {
		log.Println(err)
		resp = &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误：获取头像失败",
		}
		log.Println(avatar, resp.StatusMsg)
		return resp, nil
	}
	// 获取背景
	backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, usr.BackgroundImage)
	if err != nil {
		log.Println(err)
		resp = &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误：获取背景图失败",
		}
		log.Println(backgroundImage, resp.StatusMsg)
		return resp, nil
	}

	//返回结果
	resp = &user.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		User: &user.User{
			Id:              int64(usr.ID),
			Name:            usr.UserName,
			FollowCount:     int64(usr.FollowingCount),
			FollowerCount:   int64(usr.FollowerCount),
			IsFollow:        true,
			Avatar:          avatar,
			BackgroundImage: backgroundImage,
			Signature:       usr.Signature,
			TotalFavorited:  int64(usr.TotalFavorited),
			WorkCount:       int64(usr.WorkCount),
			FavoriteCount:   int64(usr.FavoriteCount),
		},
	}
	return resp, nil
}
