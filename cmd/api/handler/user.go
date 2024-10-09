package handler

import (
	"context"
	"net/http"
	"strconv"
	"tiktok/internal/response"
	"tiktok/kitex/kitex_gen/user"
	"tiktok/kitex/kitex_gen/user/userservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	userClient userservice.Client
)

func InitUser() {
	etcd_address := Config.GetString("etcd.address")
	// 服务发现
	resolver, err := etcd.NewEtcdResolver([]string{etcd_address})
	if err != nil {
		logger.Fatal(err)
	}
	c, err := userservice.NewClient("user", client.WithResolver(resolver))
	if err != nil {
		logger.Fatal(err)
	}
	userClient = c
	logger.Println("init user client")
}

// Register 注册
func Register(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	logger.Println("register", username)
	//校验参数
	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusBadRequest, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码不能为空",
			},
		})
		return
	}
	if len(username) > 32 || len(password) > 32 {
		c.JSON(http.StatusBadRequest, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码长度不能大于32个字符",
			},
		})
		return
	}
	//调用kitex/kitex_gen
	req := &user.UserRegisterRequest{
		Username: username,
		Password: password,
	}
	res, _ := userClient.Register(ctx, req)
	if res == nil {
		logger.Println("user无服务")
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "无服务",
			},
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.Register{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		UserID: res.UserId,
		Token:  res.Token,
	})
}

// Login 登录
func Login(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	logger.Println("login", username)
	//校验参数
	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusBadRequest, response.Login{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码不能为空",
			},
		})
		return
	}
	//调用kitex/kitex_gen
	req := &user.UserLoginRequest{
		Username: username,
		Password: password,
	}
	res, _ := userClient.Login(ctx, req)
	if res == nil {
		logger.Println("user无服务")
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "无服务",
			},
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Login{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.Login{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		UserID: res.UserId,
		Token:  res.Token,
	})
}

// UserInfo 用户信息
func UserInfo(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	logger.Println("userInfo id", userId)
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.UserInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
			User: nil,
		})
		return
	}
	//调用kitex/kitex_gen
	req := &user.UserInfoRequest{
		UserId: id,
	}
	res, _ := userClient.UserInfo(ctx, req)
	if res == nil {
		logger.Println("user无服务")
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "无服务",
			},
		})
		return
	}
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.UserInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
			User: nil,
		})
		return
	}
	c.JSON(http.StatusOK, response.UserInfo{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		User: res.User,
	})
}
