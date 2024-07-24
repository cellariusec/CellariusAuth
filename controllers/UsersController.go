package controllers

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"cellariusauth/util"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



func Signup(c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Usertype,PaymentPlan")

	var body struct {
        ID uint
		Email    string
		Password string
		Cedula string
		Usertype string
		PaymentPlan string
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

	user := models.User{Email: body.Email, Password: string(hash), Cedula: body.Cedula, Usertype: c.GetHeader("Usertype")}
	result := initializer.DB.Create(&user)
	fmt.Println("resultado: ", result)

	fmt.Print("id: ",user.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while creating user"})
		return
	}
	signupResponse := models.SignupResponse{ID: user.ID}
	fmt.Print(user.ID)
	fmt.Println("signupResponse: ", signupResponse)
	c.JSON(http.StatusOK, signupResponse)
}



func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
		OTPCode  string
		UserType string
		Cedula   string
		Usertype string
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

userType := user.Usertype

 sessionjwt,err := util.GenerateJWTs(c,user.Email,strconv.Itoa(int(user.ID)),userType)
 fmt.Println("userid: ",user.ID)
 if err != nil{
	 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	 return
 }
 RefreshToken,err := util.GenerateRefreshToken(c,user.ID)
 if err != nil{
	 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	 return
	  }



	  cookie := &http.Cookie{
        Name:     "token",
        Value:    sessionjwt,
        Path:     "/",
        Expires:  time.Now().Add(24 * time.Hour), 
        HttpOnly: true,
    }
    c.SetCookie("token", sessionjwt, 3600*24, "", "", false, true)
    http.SetCookie(c.Writer, cookie)

	refreshCookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    RefreshToken,
        Path:     "/",
        Expires:  time.Now().Add(7 * (24 * time.Hour)), 
        HttpOnly: true,
    }
    c.SetCookie("refresh_token", RefreshToken, 3600*24, "", "", false, true)
    http.SetCookie(c.Writer, refreshCookie)




	fmt.Println(string(sessionjwt))

	c.JSON(http.StatusOK, gin.H{"token": sessionjwt, "refresh_token": RefreshToken})
}



func Validate(c *gin.Context) {
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
		userID := claims["sub"].(float64)
		username := claims["email"].(string)

		c.JSON(http.StatusOK, gin.H{"message": "access granted", "userID": userID, "username": username})

	} else {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authenticate token"})

	}
}


func Logout(c *gin.Context) {
    var requestBody struct {
        Token string `json:"token"`
    }
	dsn := "postgres://yjfgskzw:PRSNUNR2F8X8InPBIra5yi5xqozxMtx0@kala.db.elephantsql.com/yjfgskzw"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Writer.Write([]byte("Error connecting to the database"))
		return
	}
    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind request body"})
        return
    }

	expires_at := util.DecodeJWT(c,requestBody.Token)
	expiresAt := time.Unix(expires_at, 0)

    revokedToken := models.RevokedToken{Token: requestBody.Token,Expiry:expiresAt}
    result := db.Create(&revokedToken)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
        return
    }

    c.SetCookie("token", "", -1, "", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})

}

//c.JSON(http.StatusOK, gin.H{"token": tokenString})
/*
curl -X POST -H "Content-Type: application/json" -d '{"Email":"example@example.com","Password":"password123"}' http://localhost:6050/signup
*/

/*
curl -X POST -H "Content-Type: application/json" -d '{"Email":"example@example.com","Password":"password123"}' http://localhost:8080/login
*/



/*
curl -X POST -H "Content-Type: application/json" -d '{"Email":"example@example.com","Password":"password123"}' https://cellariusauth-production.up.railway.app/login

curl -X POST -H "Content-Type: application/json" -d '{"Email":"jose@google.com","Password":"password1234"}' http://localhost:6050/login
*/

