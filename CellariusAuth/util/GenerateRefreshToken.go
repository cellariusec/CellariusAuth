package util

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



func GenerateRefreshToken(c *gin.Context,userID uint) (string, error) {
var existingToken models.RefreshToken
err := initializer.DB.Where("user_id=?", userID).First(&existingToken).Error

if err == nil {
    return existingToken.Token, nil

}

    refreshToken := uuid.NewString()
    expiresAt := time.Now().Add(time.Hour * 24 * 30)
	fmt.Println(refreshToken)
    err = initializer.DB.Create(&models.RefreshToken{
        Token:  refreshToken,
        UserID: userID,
        ExpiresAt: expiresAt,
    }).Error

    if err != nil {
        return "", err
    }
c.JSON(http.StatusOK, gin.H{"token": refreshToken})
    return refreshToken, nil
}
