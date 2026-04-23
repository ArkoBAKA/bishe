package model

import "time"

type Like struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"uniqueIndex:uk_likes_user_target;index;not null" json:"userId"`
	TargetType string    `gorm:"size:16;uniqueIndex:uk_likes_user_target;index;not null" json:"targetType"`
	TargetID   uint64    `gorm:"uniqueIndex:uk_likes_user_target;index;not null" json:"targetId"`
	CreatedAt  time.Time `json:"createdAt"`
}
