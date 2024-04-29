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
config.AddAllowMethods("GET", "POST", "PUT", "DELETE","OPTIONS")
config.AddAllowHeaders("Authorization", "Content-Type")
r.Use(cors.New(config))

r.OPTIONS("/*path",func(c *gin.Context){
	if c.Request.Method == "OPTIONS" {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS") // Update to include OPTIONS
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Status(204)
		return
	}
	c.Next()
})
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

