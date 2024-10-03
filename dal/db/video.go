package db

import "gorm.io/gorm"

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