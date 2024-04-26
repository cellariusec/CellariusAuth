package controllers

import (
	"fmt"
	"time"
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"cellariusauth/util"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	
)



func Signup(c *gin.Context) {

	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while hashing password",
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializer.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}



func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
		OTPCode  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form"})
		return
	}

	var user models.User
	result := initializer.DB.Where("email = ?", body.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usuario no existe!"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

 sessionjwt,err := util.GenerateJWTs(c,user.Email,string(rune(user.ID)),"admin")
 if err != nil{
	 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	 return
 }
 util.GenerateRefreshToken(c,user.ID)



	cookie := &http.Cookie{
		Name:     "token",
		Value:    sessionjwt,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	c.SetCookie("RefreshToken", sessionjwt, 3600*24*30, "", "", false, true)
	http.SetCookie(c.Writer, cookie)


	fmt.Println(string(sessionjwt))
	//fmt.Println(string(refreshToken))

	c.JSON(http.StatusOK, gin.H{"token": sessionjwt})
}



func Validate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}


	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Metodo de firma no esperado:%v", token.Header["alg"])

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
		userID := claims["sub"].(float64)
		username := claims["email"].(string)

		c.JSON(http.StatusOK, gin.H{"message": "access granted", "userID": userID, "username": username})

	} else {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authenticate token"})

	}
}
//c.JSON(http.StatusOK, gin.H{"token": tokenString})
/*
curl -X POST -H "Content-Type: application/json" -d '{"Email":"example@example.com","Password":"password123"}' http://localhost:8080/signup
*/

/*
curl -X POST -H "Content-Type: application/json" -d '{"Email":"example@example.com","Password":"password123"}' http://localhost:8080/login
*/



/*
curl -X POST -H "Content-Type: application/json" -d '{"Email":"example@example.com","Password":"password123"}' https://independent-sparkle-production.up.railway.app/login

*/