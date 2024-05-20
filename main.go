package main

import (
	"cellariusauth/controllers"
	initializer "cellariusauth/initializers"
	"cellariusauth/util"
    "cellariusauth/middleware"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
    "os"
)

func init() {
    os.Setenv("DB_CONNECTION_STRING", "postgres://yjfgskzw:PRSNUNR2F8X8InPBIra5yi5xqozxMtx0@kala.db.elephantsql.com/yjfgskzw")
	os.Setenv("ISSUER", "http://localhost:8080")
	os.Setenv("SECRET","secret")
	os.Setenv("JWT_SECRET","secret")
	os.Setenv("AUDIENCE","http://localhost:5000")
    initializer.LoadEnvVariables()
    initializer.ConnectToDb()
    initializer.SyncDatabase()
}

func main() {
  
    blacklistDSN := "postgres://yjfgskzw:PRSNUNR2F8X8InPBIra5yi5xqozxMtx0@kala.db.elephantsql.com/yjfgskzw"
    blacklistDB, err := gorm.Open(postgres.Open(blacklistDSN), &gorm.Config{})
    if err != nil {
        // Handle error
        return
    }
    r := gin.Default()

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000", "https://verbose-orbit-ppxqpr7g9q6f6g4.github.dev","http://localhost:8080"} 
    config.AllowCredentials = true
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Cookie"} 

    r.Use(cors.New(config))

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    r.POST("/signup", middleware.ValidateMiddleware, controllers.Signup)
    r.POST("/login", controllers.Login)
    r.GET("/validate", controllers.Validate)
    r.POST("/refresh-token", controllers.RefreshToken)
    r.POST("/logout", controllers.Logout)

    ticker := time.NewTicker(1*time.Hour)

    go func(){
        for range ticker.C{
            util.DeleteExpiredTokens(blacklistDB)
        }
    }()

    r.Run()
}