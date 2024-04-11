package util

import (
    "fmt"
   // "time"
   "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
)

func GenerateJWTs(c*gin.Context)(string, error) {
    // Create a new token object with the desired claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub":   "1234567890",
        "name":  "John Doe",
        "admin": true,
        "exp":   1712704106,
        "iss":   "https://dev-af7assovbom6a28o.us.auth0.com/",
        "aud":   "http://localhost:5000/jwt",
    })

    // Set the custom header fields
    token.Header["alg"] = "HS256"
    token.Header["typ"] = "JWT"
    token.Header["KeyId"] = "rQuFYlFDGSAGXq3JFFvwiLHygSsN51PH"

    // Sign the token with a secret key
    secret := []byte("rQuFYlFDGSAGXq3JFFvwiLHygSsN51PH")
    tokenString, err := token.SignedString(secret)
    if err != nil {
        panic(err)
    }

    // Print the token string
    fmt.Println(tokenString)
	c.JSON(200, gin.H{"token": tokenString})
	return tokenString, nil
}


/*
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsIktleUlkIjoiclF1RllsRkRHU0FHWHEzSkZGdndpTEh5Z1NzTjUxUEgifQ.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNzEyNzA0MTA2LCJuYW1lIjoiSm9obiBEb2UiLCJzdWIiOiIxMjM0NTY3ODkwIiwiaXNzIjoiaHR0cHM6Ly9kZXYtYWY3YXNzb3Zib202YTI4by51cy5hdXRoMC5jb20vIiwiYXVkIjoiaHR0cDovL2xvY2FsaG9zdDo1MDAwL2p3dCJ9.ZoNXlMbq2ZIlbSyeaHvq2svQ0agQzudEF2Qwir9ebRI
*/