package controllers

import (
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"cellariusauth/util"
	"encoding/json"
	"time"

	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	//"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Setenv("DB_CONNECTION_STRING", "postgres://yjfgskzw:PRSNUNR2F8X8InPBIra5yi5xqozxMtx0@kala.db.elephantsql.com/yjfgskzw")
	os.Setenv("ISSUER", "http://localhost:8080")
	os.Setenv("SECRET","secret")
	os.Setenv("JWT_SECRET","secret" )
	os.Setenv("AUDIENCE","http://localhost:5000")
	initializer.LoadEnvVariables()
	initializer.ConnectToDb()
	initializer.SyncDatabase()

	m.Run()
}


func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login", Login)
	return r
}
// Pruebas creacion de cuentas y login 
func TestLoginSuccess(t *testing.T) {
    // Generate a unique email address for testing
    email := fmt.Sprintf("josenaranjo%d@xmail.com", time.Now().Unix())

    // Create the user for testing
    user := createUser(t, email, "password123")

    // Perform login using the generated email and password
    reqBody := fmt.Sprintf(`{"Email": "%s", "Password": "password123"}`, email)
    req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r := setupRouter()
    r.ServeHTTP(w, req)

    // Check if login was successful
    var response struct {
        AccessToken string `json:"access_token"`
    }
    err := json.NewDecoder(w.Body).Decode(&response)
    assert.NoError(t, err)
 

    // Delete the user after the test
    defer deleteUser(t, user)
}


func createUser(t *testing.T, emailPrefix, password string) models.User {
    tx := initializer.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

    // Generate a unique email address for testing
    email := fmt.Sprintf("%s-%d@xmail.com", emailPrefix, time.Now().Unix())
    user := models.User{Email: email, Password: password}

    if err := tx.Create(&user).Error; err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }

    assert.NotZero(t, user.ID, "User ID should not be zero after creation")
    return user
}

func deleteUser(t *testing.T, user models.User) {
    if err := initializer.DB.Unscoped().Delete(&user).Error; err != nil {
        t.Errorf("Failed to delete user: %v", err)
    }
}

// Pruebas con credenciales invalidas

func TestLoginInvalidCredentials(t *testing.T) {
	 
	r := setupRouter()

 
	reqBody := `{"Email": "nonexistent@example.com", "Password": "wrongpassword"}`
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

 
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error": "Usuario no existe!"`)
}

// pruebas Logout 

func TestLogout(t *testing.T) {
    // Start a transaction
    tx := initializer.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            tx.Rollback() // Rollback the transaction after the test
        }
    }()

    // Create a user with a unique email address
    email := "test@example.com"
    user := models.User{Email: email, Password: "testpassword"}
    if err := tx.Create(&user).Error; err != nil {
        t.Errorf("Failed to create user: %v", err)
        return
    }

    // Create a gin context for testing
    r := setupRouter()
    reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email, user.Password)
    req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    c, _ := gin.CreateTestContext(w)

    token, err := util.GenerateJWTs(c, user.Email, string(rune(user.ID)), "admin")
    if err != nil {
        t.Errorf("Failed to generate token: %v", err)
        return
    }

    reqBody = fmt.Sprintf(`{"token": "%s"}`, token)
    req, _ = http.NewRequest(http.MethodPost, "/logout", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")

    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code, "Logout should succeed with valid token")

    var revokedToken models.RevokedToken
    if err := tx.Where("token = ?", token).First(&revokedToken).Error; err != nil {
        t.Errorf("Failed to find revoked token: %v", err)
    }

    if err := tx.Unscoped().Delete(&user).Error; err != nil {
        t.Errorf("Failed to delete user: %v", err)
    }
    if err := tx.Unscoped().Delete(&revokedToken).Error; err != nil {
        t.Errorf("Failed to delete revoked token: %v", err)
    }
}


