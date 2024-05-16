package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


func ValidateMiddleware (c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}


	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
		Usertype := claims["usertype"].(string)
        requiredUsertype := c.GetHeader("Usertype")
        if requiredUsertype == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Usertype header is required"})
            return
        }

		fmt.Println("usertype:",requiredUsertype)

		if !IsAuthorizedToCreateUser(Usertype, requiredUsertype) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to create this user"})
			return
		}
        c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authenticate token"})

	}
}

func IsAuthorizedToCreateUser(requestingUsertype, requiredUsertype string) bool {
	switch requestingUsertype {
	case "superadmin":
		return true

	case "admin":
		return requiredUsertype == "store" || requiredUsertype == "user"

    case "store":
		return requiredUsertype == "user"
	case"user":
		return false

	default:
		return false
	}


}