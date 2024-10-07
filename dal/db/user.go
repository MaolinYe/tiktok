package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// 用户数据
type User struct {
	gorm.Model
	UserName        string  `gorm:"index:idx_username,unique;type:varchar(40);not null" json:"name,omitempty"`
	Password        string  `gorm:"type:varchar(256);not null" json:"password,omitempty"`
	FavoriteVideos  []Video `gorm:"many2many:user_favorite_videos" json:"favorite_videos,omitempty"`
	FollowingCount  uint    `gorm:"default:0;not null" json:"follow_count,omitempty"`                                                           // 关注总数
	FollowerCount   uint    `gorm:"default:0;not null" json:"follower_count,omitempty"`                                                         // 粉丝总数
	Avatar          string  `gorm:"type:varchar(256);default:default_avatar.jpg" json:"avatar,omitempty"`                                                                  // 用户头像
	BackgroundImage string  `gorm:"type:varchar(256);default:default_background.jpg" json:"background_image,omitempty"` // 用户个人页顶部大图
	WorkCount       uint    `gorm:"default:0;not null" json:"work_count,omitempty"`                                                             // 作品数
	FavoriteCount   uint    `gorm:"default:0;not null" json:"favorite_count,omitempty"`                                                         // 喜欢数
	TotalFavorited  uint    `gorm:"default:0;not null" json:"total_favorited,omitempty"`                                                        // 获赞总量
	Signature       string  `gorm:"type:varchar(256)" json:"signature,omitempty"`                                                               // 个人简介
}

// 根据用户id获取用户数据
func GetUserByID(ctx context.Context, userID int64) (*User, error) {
	user := new(User)
	if err := db.Clauses(dbresolver.Read).WithContext(ctx).
		First(&user, userID).Error; err == nil {
		return user, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 新增用户
func CreateUser(ctx context.Context, user *User) error {
	err := db.Clauses(dbresolver.Write).WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(user).Error; err != nil {
				return err
			}
			return nil
		})
	return err
}

// 根据用户名获取密码、ID
func GetUserByUsername(ctx context.Context, userName string) (*User, error) {
	user := new(User)
	if err := db.Clauses(dbresolver.Read).WithContext(ctx).
		Select("id, password").Where("user_name = ?", userName).
		First(&user).Error; err == nil {
		return user, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
