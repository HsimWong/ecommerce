package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;size:50"`
	Email     string `gorm:"uniqueIndex;size:100"`
	Phone     string `gorm:"uniqueIndex;size:15"`
	Password  string `gorm:"size:100"` // 存储加密后的密码
	Gender    string `gorm:"size:16"`
	Salt      string `gorm:"size:30"`   // 密码盐值
	Status    int    `gorm:"default:1"` // 1-正常 2-禁用
	CreatedAt time.Time
	UpdatedAt time.Time
}
