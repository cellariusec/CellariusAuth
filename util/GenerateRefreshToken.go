package util

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



func GenerateRefreshToken(c *gin.Context, userID uint) (string, error) {
    var existingToken models.RefreshToken

    err := initializer.DB.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Order("expires_at DESC").First(&existingToken).Error
    if err == nil {
   
        return existingToken.Token, nil
    }
    if err == gorm.ErrRecordNotFound {
        refreshToken := uuid.NewString()
        expiresAt := time.Now().Add(time.Hour * 24 * 30)

        err = initializer.DB.Create(&models.RefreshToken{
            Token:     refreshToken,
            UserID:    userID,
            ExpiresAt: expiresAt,
        }).Error

        if err != nil {
            return "", err
        }

        c.JSON(http.StatusOK, gin.H{"token": refreshToken})
        return refreshToken, nil
    }



    initializer.DB.Where("user_id = ? AND expires_at < ?", userID, time.Now()).Delete(&models.RefreshToken{})

 
    refreshToken := uuid.NewString()
    expiresAt := time.Now().Add(time.Hour * 24 * 30)

    err = initializer.DB.Create(&models.RefreshToken{
        Token:     refreshToken,
        UserID:    userID,
        ExpiresAt: expiresAt,
    }).Error

    if err != nil {
        return "", err
    }

    c.JSON(http.StatusOK, gin.H{"token": refreshToken})
    return refreshToken, nil
}