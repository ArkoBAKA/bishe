package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID              uint64         `gorm:"primaryKey" json:"commentId"`
	PostID          uint64         `gorm:"index;not null" json:"postId"`
	AuthorID        uint64         `gorm:"index;not null" json:"authorId"`
	ParentCommentID uint64         `gorm:"index;not null;default:0" json:"parentCommentId"`
	Content         string         `gorm:"size:1000;not null" json:"content"`
	LikeCount       int64          `gorm:"not null;default:0" json:"likeCount"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
