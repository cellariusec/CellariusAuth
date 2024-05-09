package models

import (
	"time"
	"gorm.io/gorm"

)

type RevokedToken struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"type:varchar(2000);unique;not null"`
	Expiry time.Time    `gorm:not null`
}
