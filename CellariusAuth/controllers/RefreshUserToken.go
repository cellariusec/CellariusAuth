package controllers

import (
    "fmt"
    initializer "cellariusauth/initializers"
    "cellariusauth/models"
    "net/http"
    "github.com/gin-gonic/gin"
    "cellariusauth/util"
)

func RefreshToken(c *gin.Context) {
    var refreshTokenData struct {
        RefreshToken string `json:"refresh_token"`
    }

    if err := c.BindJSON(&refreshTokenData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        fmt.Println(err)
        return
    }

    refreshTokenString := refreshTokenData.RefreshToken

    var refreshToken models.RefreshToken

    if err := initializer.DB.Where("token = ?", refreshTokenString).First(&refreshToken).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
        fmt.Println(err)
        return
    }

    var user models.User

    if err := initializer.DB.First(&user, refreshToken.UserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        fmt.Println(err)
        return
    }

    tokenString, err := util.GenerateJWTs(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        fmt.Println(err)
        return
    }
    fmt.Println("new sessio", tokenString)
    //c.JSON(http.StatusOK, gin.H{"token": tokenString})


}