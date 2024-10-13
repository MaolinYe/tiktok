package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Favor struct {
	Video   Video `gorm:"foreignkey:VideoID;" json:"video,omitempty"`
	VideoID uint  `gorm:"index:idx_videoid;not null" json:"video_id"`
	User    User  `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID  uint  `gorm:"index:idx_userid;not null" json:"user_id"`
}

// 新增点赞记录，和视频相关的操作交给redis处理
func CreateFavor(ctx context.Context, favor *Favor) error {
	err := db.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 新增点赞记录
		if err := tx.Create(favor).Error; err != nil {
			return err
		}
		// 增加用户点赞数量
		if err := tx.Model(&User{}).Where("id = ?", favor.UserID).
			Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteFavor 删除点赞记录
func DeleteFavor(ctx context.Context, favor *Favor) error {
	err := db.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除点赞记录
		if err := tx.Where("user_id = ? AND video_id = ?", favor.UserID, favor.VideoID).
			Delete(&Favor{}).Error; err != nil {
			return err
		}
		// 减少用户点赞数量
		if err := tx.Model(&User{}).Where("id = ?", favor.UserID).
			Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// 更新视频点赞数量
func UpdateFavorCount(ctx context.Context, videoID int64, favorCount int64) error {
	err := db.Clauses(dbresolver.Write).WithContext(ctx).
		First(&Video{}, videoID).Update("favorite_count", favorCount).Error
	return err
}

// 获取视频点赞数量
func GetFavorCount(ctx context.Context, videoID int64) (uint, error) {
	video := new(Video)
	err := db.Clauses(dbresolver.Read).WithContext(ctx).Select("favorite_count").First(&video, videoID).Error
	return video.FavoriteCount, err
}

// 获取用户点赞的视频列表
func GetFavorList(ctx context.Context, userID int64) (list []*Video, err error) {
	// 查询该用户所有的 Favor 记录并预加载视频信息
	var favors []Favor
	err = db.Clauses(dbresolver.Read).WithContext(ctx).Preload("Video").Where("user_id = ?", userID).Find(&favors).Error
	if err != nil {
		return nil, err
	}
	// 提取视频列表
	for _, favor := range favors {
		list = append(list, &favor.Video)
	}
	return
}

// 查看用户是否点赞该视频
func IsFavorite(ctx context.Context, videoID uint, userID int64) (bool, error) {
	var favor Favor
	err := db.Clauses(dbresolver.Read).WithContext(ctx).Where("user_id = ? and video_id = ?", userID, videoID).First(&favor).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
