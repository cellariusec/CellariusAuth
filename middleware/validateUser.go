package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ValidateEmailMiddleware (c *gin.Context) {

	
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}


	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo de firma no esperado:%v", token.Header["alg"])

		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authenticate token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("Claims: %+v", claims)
		Email := claims["username"].(string)
		c.Set("email", Email) 

	
        c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authenticate token"})

	}
}

