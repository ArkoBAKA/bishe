package model

import (
	"time"

	"gorm.io/gorm"
)

type UploadFile struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	Bucket       string         `gorm:"size:32;index;not null" json:"bucket"`
	Scene        string         `gorm:"size:32;index;not null;default:''" json:"scene"`
	OriginalName string         `gorm:"size:255;not null" json:"originalName"`
	StoredName   string         `gorm:"size:80;uniqueIndex;not null" json:"storedName"`
	Ext          string         `gorm:"size:16;index;not null" json:"ext"`
	Size         int64          `gorm:"not null" json:"size"`
	MimeType     string         `gorm:"size:128;not null" json:"mimeType"`
	RelPath      string         `gorm:"size:512;not null" json:"relPath"`
	URL          string         `gorm:"size:512;not null" json:"url"`
	UploaderID   uint64         `gorm:"index;not null" json:"uploaderId"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
