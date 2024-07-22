package controllers

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"fmt"
	"net/http"
	"os"
	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

func ResetPassword(c *gin.Context) {
    var body struct {
        CurrentPassword string `json:"current_password"`
        NewPassword     string `json:"new_password"`
    }
    if err := c.BindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    tokenString := c.GetHeader("Authorization")
    if tokenString == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
        return
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("SECRET")), nil
    })

    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
        return
    }

    email, ok := claims["username"].(string)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
        return
    }

    var user models.User
    result := initializer.DB.Where("LOWER(email) = LOWER(?)", email).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Verify current password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.CurrentPassword))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password"})
        return
    }

    // Hash and update new password
    hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user.Password = string(hash)
    if err := initializer.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}