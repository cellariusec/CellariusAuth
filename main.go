package main

import (
	"cellariusauth/controllers"
	initializer "cellariusauth/initializers"
	"cellariusauth/middleware"
	"cellariusauth/util"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	os.Setenv("DB_CONNECTION_STRING", "postgres://avnadmin:AVNS_LNgImquHJXNIMn4aMTt@147.182.201.146:18022/defaultdb?sslmode=require")
	os.Setenv("ISSUER", "http://localhost:8080")
	os.Setenv("SECRET", "secret")
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("AUDIENCE", "http://localhost:5000")
	initializer.LoadEnvVariables()
	initializer.ConnectToDb()
	initializer.SyncDatabase()
}

func main() {

	blacklistDSN := "postgres://avnadmin:AVNS_LNgImquHJXNIMn4aMTt@147.182.201.146:18022/defaultdb?sslmode=require"
	blacklistDB, err := gorm.Open(postgres.Open(blacklistDSN), &gorm.Config{})
	if err != nil {
		// Handle error
		return
	}
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000",
		"https://cellariusec-cellarius-web-store.vercel.app",
		"https://cellariusec-cellarius-web-store-icu5c4pzw-cellarius-projects.vercel.app",
		"https://cellariusec-cellarius-web-store-git-main-cellarius-projects.vercel.app",
		"http://localhost:8080"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Cookie", "Usertype"}

	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.OPTIONS("/signup", func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		c.AbortWithStatus(http.StatusOK)
	})
	//middleware.ValidateMiddleware
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", controllers.Validate)
	r.POST("/refresh-token", controllers.RefreshToken)
	r.POST("/logout", controllers.Logout)
	r.POST("/reset_password", middleware.ValidateEmailMiddleware, controllers.ResetPassword)
	r.POST("/get_device_token", controllers.GetDeviceToken)

	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			util.DeleteExpiredTokens(blacklistDB)
		}
	}()

	r.Run(":5050")
}
