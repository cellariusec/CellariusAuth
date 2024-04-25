package main

import (
	"bytes"
	"cellariusauth/controllers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type signupRequestBody struct {
    Email    string `json:"Email"`
    Password string `json:"Password"`
}
// testing del signup(creacion de usuario)
func TestSignup(t *testing.T) {
router := gin.Default()
router.POST("/signup", controllers.Signup)
requestBody := signupRequestBody{
	Email:    "jnarcursos1234@gmail.com",
	Password: "password123",
}
jsonBody, _ := json.Marshal(requestBody)

w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
req.Header.Set("Content-Type", "application/json")
router.ServeHTTP(w, req)
assert.Equal(t, 200, w.Code)
//assert.Equal(t, "User created successfully", w.Body.String())
}
// testing de invalid form data 

func TestFormData(t* testing.T){
	router := gin.Default()
	router.POST("/login", controllers.Login)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid form", response["error"])
}






