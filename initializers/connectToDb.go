package initializer

import (
	"fmt"
	"net/url"
	"cellariusauth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"os"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error

	encodedPassword := url.QueryEscape("AVNS_LNgImquHJXNIMn4aMTt")

	
	dsn := fmt.Sprintf("postgres://avnadmin:%s@actixwebpostgres-udla-54df.aivencloud.com:18022/defaultdb?sslmode=require", encodedPassword)
	//dsn := "postgres://avnadmin:AVNS_LNgImquHJXNIMn4aMTt@actixwebpostgres-udla-54df.aivencloud.com:18022/defaultdb?sslmode=require"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.RevokedToken{})
	DB.AutoMigrate(&models.RefreshToken{})

	if err != nil {
		panic("Failed to connect to database!")
	} else {
		fmt.Println("Connected to database")
	}
}
