package model

import (
	"time"

	"gorm.io/gorm"
)

type Forum struct {
	ID          uint64         `gorm:"primaryKey" json:"forumId"`
	Name        string         `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Description string         `gorm:"size:255;not null;default:''" json:"description"`
	CoverURL    string         `gorm:"size:512;not null;default:''" json:"coverUrl"`
	OwnerID     uint64         `gorm:"index;not null" json:"ownerId"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
