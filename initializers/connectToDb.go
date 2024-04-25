package initializer

import(
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
	//"os"
)

var DB *gorm.DB
func ConnectToDb(){
	var err error
	dsn := "postgres://yjfgskzw:PRSNUNR2F8X8InPBIra5yi5xqozxMtx0@kala.db.elephantsql.com/yjfgskzw"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		panic("Failed to connect to database!")
	}else{
		fmt.Println("Connected to database")
	}
}