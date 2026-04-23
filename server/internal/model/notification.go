package model

import (
	"time"
)

type Notification struct {
	ID        uint64     `gorm:"primaryKey" json:"notificationId"`
	UserID    uint64     `gorm:"index;not null" json:"userId"`
	Type      string     `gorm:"size:16;index;not null" json:"type"`
	Title     string     `gorm:"size:64;not null;default:''" json:"title"`
	Content   string     `gorm:"size:255;not null;default:''" json:"content"`
	DataJSON  string     `gorm:"size:2000;not null;default:''" json:"-"`
	IsRead    bool       `gorm:"index;not null;default:false" json:"isRead"`
	ReadAt    *time.Time `json:"readAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
