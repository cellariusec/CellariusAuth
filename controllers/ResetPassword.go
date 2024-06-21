package controllers

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ResetPassword(c *gin.Context) {
	var body struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BindJSON(&body); err != nil {
		fmt.Println("Error in binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tokenstring := c.GetHeader("Authorization")
	if tokenstring == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error al autenticar!"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenEmail := claims["username"].(string)

		var user models.User
		fmt.Println("Searching for user with email:", body.Email)
		fmt.Println(user)

		// Convert both email addresses to lower case for case-insensitive comparison
		body.Email = strings.ToLower(body.Email)
		tokenEmail = strings.ToLower(tokenEmail)

		result := initializer.DB.Where("LOWER(email) = ?", body.Email).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				fmt.Println("Error querying database:", result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		if user.Email != tokenEmail {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		fmt.Println("Current password:", body.Password)
		fmt.Println("New password:", body.NewPassword)

		hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.Password = string(hash)
		if err := initializer.DB.Save(&user).Error; err != nil {
			fmt.Println("Error saving user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authenticate token"})
	}
}