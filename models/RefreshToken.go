package models

import (
	"time"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	Token  string `gorm:"unique;column:token"`
	UserID uint
	ExpiresAt time.Time
}
