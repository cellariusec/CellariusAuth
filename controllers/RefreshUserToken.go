package controllers

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"cellariusauth/util"
	"fmt"
	"net/http"
	//"time"
	"github.com/gin-gonic/gin"
)

func RefreshToken(c *gin.Context) {
    //refreshTokenString := c.GetHeader("RefreshToken")
    var refreshTokenData struct {
        RefreshToken string `json:"refresh_token"`
    }

    if err := c.BindJSON(&refreshTokenData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        fmt.Println(err)
        return
    }

    refreshTokenString := refreshTokenData.RefreshToken

    var user models.User

    var refreshToken models.RefreshToken

    if err := initializer.DB.Where("token = ?", refreshTokenString).First(&refreshToken).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
        fmt.Println(err)
        return
    }
/*

    if refreshToken.ExpiresAt.Before(time.Now()){
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
        fmt.Println("Refresh token expired")
        return
    }



    if err := initializer.DB.First(&user, refreshToken.UserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        fmt.Println(err)
        return
    }
    */

    tokenString, err := util.GenerateJWTs(c, user.Email, string(user.ID), "admin")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        fmt.Println(err)
        return
    }
    fmt.Println("new session", tokenString)
    c.JSON(http.StatusOK, gin.H{"token": tokenString})


}