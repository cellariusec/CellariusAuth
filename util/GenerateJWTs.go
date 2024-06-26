package util

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateJWTs(c*gin.Context,username,userid,usertype string)(string, error) {
    expirationTime := time.Now().Add(168 * time.Hour)
    fmt.Println("userid generate jwt: ",userid)
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

    cookieName := "jwt_token"
    cookieMaxAge := int(48 * time.Hour/time.Second)
    cookieDomain := "localhost"
    cookiePath := "/"
    cookieSecure := true
    cookieHttpOnly := true
    //cookieSameSite := "strict"

    cookie := &http.Cookie{
        Name: cookieName,
        Value: tokenString,
        MaxAge: cookieMaxAge,
        Domain: cookieDomain,
        Path: cookiePath,
        Secure: cookieSecure,
        HttpOnly: cookieHttpOnly,
        //SameSite: http.SameSiteMode(cookieSameSite),
    }
    c.SetCookie(cookieName, tokenString, cookieMaxAge, cookiePath, cookieDomain, cookieSecure, cookieHttpOnly)
    c.Writer.Header().Set("Set-Cookie", cookie.String())
    return tokenString, nil
}


/*
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsIktleUlkIjoiclF1RllsRkRHU0FHWHEzSkZGdndpTEh5Z1NzTjUxUEgifQ.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNzEyNzA0MTA2LCJuYW1lIjoiSm9obiBEb2UiLCJzdWIiOiIxMjM0NTY3ODkwIiwiaXNzIjoiaHR0cHM6Ly9kZXYtYWY3YXNzb3Zib202YTI4by51cy5hdXRoMC5jb20vIiwiYXVkIjoiaHR0cDovL2xvY2FsaG9zdDo1MDAwL2p3dCJ9.ZoNXlMbq2ZIlbSyeaHvq2svQ0agQzudEF2Qwir9ebRI
*/