package initializer

import(
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
	"os"
)

var DB *gorm.DB
func ConnectToDb(){
	var err error
	dsn := os.Getenv("DB_CONNECTION_STRING")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		panic("Failed to connect to database!")
	}else{
		fmt.Println("Connected to database")
	}
}