package model

import (
	"time"
)

type Report struct {
	ID            uint64     `gorm:"primaryKey" json:"reportId"`
	ReporterID    uint64     `gorm:"index;not null" json:"reporterId"`
	TargetType    string     `gorm:"size:16;index;not null" json:"targetType"`
	TargetID      uint64     `gorm:"index;not null" json:"targetId"`
	Reason        string     `gorm:"size:64;not null" json:"reason"`
	Detail        string     `gorm:"size:1000;not null;default:''" json:"detail"`
	Status        string     `gorm:"size:16;index;not null;default:'pending'" json:"status"`
	ProcessedBy   uint64     `gorm:"not null;default:0" json:"processedBy"`
	ProcessedAt   *time.Time `json:"processedAt,omitempty"`
	ProcessAction string     `gorm:"size:32;not null;default:''" json:"processAction"`
	ProcessRemark string     `gorm:"size:255;not null;default:''" json:"processRemark"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
