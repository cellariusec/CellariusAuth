package controllers

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
func ResetPassword(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		NewPassword string `json:"new_password"`

	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tokenstring := c.GetHeader("Authorization")
	if tokenstring == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
        return
    }

	token,err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("SECRET")), nil
    })

	if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Error al autenticar!"})
        return
    }
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        tokenEmail := claims["email"].(string)


	var user models.User
	result  := initializer.DB.Where("email = ?", body.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	if user.Email != tokenEmail{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
fmt.Println(body.Password)
fmt.Println(body.NewPassword)
fmt.Println(body.Email)
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hash)
	initializer.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
}