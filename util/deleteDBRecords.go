package util

import (
    "cellariusauth/models"
    "gorm.io/gorm"
    "time"
)

func DeleteExpiredTokens(db *gorm.DB) {
    tx := db.Begin()

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if tx.Error != nil {
        return
    }

    if err := tx.Where("expires_at < ?", time.Now()).Delete(&models.RevokedToken{}).Error; err != nil {
        tx.Rollback()
        return
    }

    tx = tx.Commit()
    if tx.Error != nil {
        return
    }
}