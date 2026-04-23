package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID           uint64         `gorm:"primaryKey" json:"postId"`
	ForumID      uint64         `gorm:"index;not null" json:"forumId"`
	AuthorID     uint64         `gorm:"index;not null" json:"authorId"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Content      string         `gorm:"size:5000;not null" json:"content"`
	Status       string         `gorm:"size:16;index;not null;default:'pending'" json:"status"`
	LikeCount    int64          `gorm:"not null;default:0" json:"likeCount"`
	CommentCount int64          `gorm:"not null;default:0" json:"commentCount"`
	ViewCount    int64          `gorm:"not null;default:0" json:"viewCount"`
	ReviewedBy   uint64         `gorm:"not null;default:0" json:"reviewedBy"`
	ReviewedAt   *time.Time     `json:"reviewedAt,omitempty"`
	ReviewRemark string         `gorm:"size:255;not null;default:''" json:"reviewRemark"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
