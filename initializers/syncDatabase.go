package initializer

import "cellariusauth/models"

func SyncDatabase(){
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.RevokedToken{})
}