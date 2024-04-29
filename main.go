package main

import (
	"cellariusauth/controllers"
	initializer "cellariusauth/initializers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)



func init(){
	initializer.LoadEnvVariables()
	initializer.ConnectToDb()
	initializer.SyncDatabase()
}

func main(){

r := gin.Default()
config := cors.DefaultConfig()
config.AllowAllOrigins = true
config.AddAllowHeaders("Authorization", "Content-Type")
r.Use(cors.New(config))
r.GET("/ping", func(c *gin.Context){
	c.JSON(200, gin.H{
		"message": "pong",
	})

})

r.POST("/signup", controllers.Signup)
r.POST("/login", controllers.Login)
r.GET("/validate",controllers.Validate)
r.POST("/refresh-token",controllers.RefreshToken)
r.Run()
}

