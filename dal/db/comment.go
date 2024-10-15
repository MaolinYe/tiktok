package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Comment struct {
	gorm.Model
	Video   Video  `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID uint   `gorm:"index:idx_videoid;not null" json:"video_id"`
	User    User   `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID  uint   `gorm:"index:idx_userid;not null" json:"user_id"`
	Content string `gorm:"type:varchar(255);not null" json:"content"`
}

// 新增评论记录，和视频相关的操作交给redis处理
func CreateComment(ctx context.Context, comment *Comment) error {
	return db.Clauses(dbresolver.Write).WithContext(ctx).Create(comment).Error
}

// DeleteFavor 删除评论记录
func DeleteComment(ctx context.Context, commentID int64) error {
	return db.Clauses(dbresolver.Write).WithContext(ctx).Delete(&Comment{}, commentID).Error
}

// 更新视频评论数
func UpdateCommentCount(ctx context.Context, videoID int64, commentCount int64) error {
	err := db.Clauses(dbresolver.Write).WithContext(ctx).First(&Video{}, videoID).
		Update("comment_count", commentCount).Error
	return err
}

// 获取视频评论数
func GetCommentCount(ctx context.Context, videoID int64) (uint, error) {
	video := new(Video)
	err := db.Clauses(dbresolver.Read).WithContext(ctx).Select("comment_count").First(&video, videoID).Error
	return video.CommentCount, err
}

// 获取评论列表
func GetCommentList(ctx context.Context, videoID int64) (comments []*Comment, err error) {
	if err = db.Clauses(dbresolver.Read).WithContext(ctx).Where("video_id = ?", videoID).
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return
}

// 根据评论id获取评论
func GetCommentByID(ctx context.Context, commentID int64) (comment *Comment, err error) {
	err = db.Clauses(dbresolver.Read).WithContext(ctx).First(&comment, commentID).Error
	return
}
