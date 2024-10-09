package db

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Video struct {
	gorm.Model
	Author        User   `gorm:"foreignkey:AuthorID" json:"author,omitempty"`
	AuthorID      uint   `gorm:"index:idx_authorid;not null" json:"author_id,omitempty"`
	PlayUrl       string `gorm:"type:varchar(255);not null" json:"play_url,omitempty"`
	CoverUrl      string `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount uint   `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	CommentCount  uint   `gorm:"default:0;not null" json:"comment_count,omitempty"`
	Title         string `gorm:"type:varchar(50);not null" json:"title,omitempty"`
}

// 新增视频
func CreateVideo(ctx context.Context, video *Video) error {
	err := db.Clauses(dbresolver.Write).WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			// 创建视频
			if err := tx.Create(video).Error; err != nil {
				return err
			}
			// 更新用户创建视频数量
			if err := tx.Model(&User{}).Where("id = ?", video.AuthorID).
				Update("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
				return err
			}
			return nil
		})
	return err
}

// 获取用户发布的视频列表
func GetPublishList(ctx context.Context, userID int64) (list []*Video, err error) {
	if err = db.Clauses(dbresolver.Read).WithContext(ctx).Model(&Video{}).Where("author_id = ?", userID).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return
}

// 获取视频流
func GetFeed(ctx context.Context, limit int, lastestTime int64) (list []*Video, err error) {
	if lastestTime == 0 {
		lastestTime = time.Now().UnixMilli()
	}
	if err = db.Clauses(dbresolver.Read).WithContext(ctx).Model(&Video{}).
		Where("created_at < ?", time.UnixMilli(lastestTime)).Order("created_at desc").
		Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}
	return
}
