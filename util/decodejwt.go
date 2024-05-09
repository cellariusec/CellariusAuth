package util

import(
	"github.com/golang-jwt/jwt/v4"
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"log"

)

func DecodeJWT(c *gin.Context,tokenString string) int64{
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo de firma no esperado:%v", token.Header["alg"])

		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authenticate token"})
		
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("Claims: %+v", claims)
		expiryDate := int64(claims["exp"].(float64))
		return expiryDate
		
	} else {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authenticate token"})

	}
	return 0
}