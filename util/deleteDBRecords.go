package util

import (
    "cellariusauth/models"
    "gorm.io/gorm"
    "time"
)

// DeleteOldRecords deletes revoked tokens older than 7 days from the database.
// It should be called as a scheduled task or background process, not on every logout request.
func DeleteOldRecords(db *gorm.DB) {
    cutOffTime := time.Now().Add(-7 * 24 * time.Hour)

    // Start a transaction
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if tx.Error != nil {
        return
    }

    // Batch delete old records
    if err := tx.Where("created_at < ?", cutOffTime).Delete(&models.RevokedToken{}).Error; err != nil {
        tx.Rollback()
        return
    }

    // Commit the transaction
    tx = tx.Commit()
    if tx.Error != nil {
        return
    }
}