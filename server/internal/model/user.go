/*
 * @Date: 2026-04-22 17:53:42
 * @LastEditors: ArongWang 3312428832@qq.com
 * @LastEditTime: 2026-04-23 10:37:32
 * @FilePath: /server/internal/model/user.go
 * @Description:
 */
package model

import "time"

type User struct {
	ID           uint64     `gorm:"primaryKey" json:"userId"`
	Username     string     `gorm:"column:username;size:64;not null;default:''" json:"username"`
	Account      string     `gorm:"size:64;uniqueIndex;not null" json:"account"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Nickname     string     `gorm:"size:64;not null" json:"nickname"`
	AvatarURL    string     `gorm:"size:512;not null;default:''" json:"avatarUrl"`
	Bio          string     `gorm:"size:255;not null;default:''" json:"bio"`
	Role         string     `gorm:"size:16;not null;default:'user'" json:"role"`
	Status       string     `gorm:"size:16;not null;default:'normal'" json:"status"`
	BanUntil     *time.Time `json:"banUntil,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
