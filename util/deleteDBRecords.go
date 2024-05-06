package util
import(

	"cellariusauth/models"
	"gorm.io/gorm"
    "time"
)

func DeleteOldRecords(db *gorm.DB) {

    cutOffTime := time.Now().Add(-7*24*time.Hour - time.Hour)
    var recordsToDelete []models.RevokedToken
    db.Where("created_at < ?", cutOffTime).Find(&recordsToDelete)
    for _, record := range recordsToDelete {
        db.Delete(&record)
    }
}