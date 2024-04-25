package util

import (
    "fmt"
   "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "os"
    "time"
)

func GenerateJWTs(c*gin.Context,username,userid,usertype string)(string, error) {
    expirationTime := time.Now().Add(48 * time.Hour)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub":   userid,
        "username":  username,
        "userid": userid,
        "usertype": usertype,
        "exp":   expirationTime.Unix(),
        "iss":   os.Getenv("ISSUER"),
        "aud":   os.Getenv("AUDIENCE"),
    })
 
    token.Header["alg"] = "HS256"
    token.Header["typ"] = "JWT"
    token.Header["KeyId"] = os.Getenv("JWT_SECRET")
    JWT_SECRET := os.Getenv("JWT_SECRET")
    secret := []byte(JWT_SECRET)
    tokenString, err := token.SignedString(secret)
    if err != nil {
        panic(err)
    }

    fmt.Println(tokenString)
	//c.JSON(200, gin.H{"token": tokenString})
	return tokenString, nil
}


/*
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsIktleUlkIjoiclF1RllsRkRHU0FHWHEzSkZGdndpTEh5Z1NzTjUxUEgifQ.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNzEyNzA0MTA2LCJuYW1lIjoiSm9obiBEb2UiLCJzdWIiOiIxMjM0NTY3ODkwIiwiaXNzIjoiaHR0cHM6Ly9kZXYtYWY3YXNzb3Zib202YTI4by51cy5hdXRoMC5jb20vIiwiYXVkIjoiaHR0cDovL2xvY2FsaG9zdDo1MDAwL2p3dCJ9.ZoNXlMbq2ZIlbSyeaHvq2svQ0agQzudEF2Qwir9ebRI
*/