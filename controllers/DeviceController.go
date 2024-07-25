package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
	"net/http"
)

type DeviceRequest struct {
	UserID string `json:"userid"`
}

func GetDeviceToken(c *gin.Context) {
	fmt.Println("GetDeviceToken")
	expirationTime := time.Now().Add(8760 * time.Hour)

	var deviceRequest DeviceRequest
	if err := c.BindJSON(&deviceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   deviceRequest.UserID,
		"exp":   expirationTime.Unix(),
		"iss":   os.Getenv("ISSUER"),
		"aud":   os.Getenv("AUDIENCE"),
	})

	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"
	token.Header["KeyId"] = os.Getenv("DEVICE_JWT_SECRET")
	JWT_SECRET := os.Getenv("DEVICE_JWT_SECRET")
	secret := []byte(JWT_SECRET)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}

	fmt.Println(tokenString)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}