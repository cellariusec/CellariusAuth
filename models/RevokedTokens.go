package models

import (
	"gorm.io/gorm"
)

type RevokedToken struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"type:varchar(255);unique;not null"`
}
 