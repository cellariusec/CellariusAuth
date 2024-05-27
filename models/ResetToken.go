package models
import("time")
type ResetToken struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"index"`
	Token     string `gorm:"unique"`
	CreatedAt time.Time
}
