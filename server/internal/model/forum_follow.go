package model

import "time"

type ForumFollow struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"uniqueIndex:uk_forum_follows_user_forum;index;not null" json:"userId"`
	ForumID   uint64    `gorm:"uniqueIndex:uk_forum_follows_user_forum;index;not null" json:"forumId"`
	CreatedAt time.Time `json:"createdAt"`
}
