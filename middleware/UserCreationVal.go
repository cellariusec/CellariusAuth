package middleware

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
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

    var revokedToken models.RevokedToken
    result := initializer.DB.Where("token = ?", tokenString).First(&revokedToken)
    if result.Error == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
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
		var requestData struct {
            Username string `json:"username"`
            Password string `json:"password"`
            UserType string `json:"usertype"`
        }
        if err := c.ShouldBindJSON(&requestData); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

		fmt.Println("usertype:",requestData.UserType)

		if !IsAuthorizedToCreateUser(Usertype, requestData.UserType) {
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