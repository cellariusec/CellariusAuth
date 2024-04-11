package initializer

import "cellariusauth/models"

func SyncDatabase(){
	DB.AutoMigrate(&models.User{})
}