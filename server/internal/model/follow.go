package model

import (
	"time"
)

type Follow struct {
	ID         uint64     `gorm:"primaryKey" json:"followId"`
	UserID     uint64     `gorm:"uniqueIndex:uk_follows_user_target;index;not null" json:"userId"`
	TargetType string     `gorm:"size:16;uniqueIndex:uk_follows_user_target;index;not null" json:"targetType"`
	TargetID   uint64     `gorm:"uniqueIndex:uk_follows_user_target;index;not null" json:"targetId"`
	Status     string     `gorm:"size:16;index;not null;default:'active'" json:"status"`
	CanceledAt *time.Time `json:"canceledAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}
