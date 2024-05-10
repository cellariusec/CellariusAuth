package controllers

import (
	"bytes"
	initializer "cellariusauth/initializers"
	"cellariusauth/models"
	"cellariusauth/util"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"

	//"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
 
    email := fmt.Sprintf("josenaranjo%d@xmail.com", time.Now().Unix())
 
    user := createUser(t, email, "password123")

 
    reqBody := fmt.Sprintf(`{"Email": "%s", "Password": "password123"}`, email)
    req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r := setupRouter()
    r.ServeHTTP(w, req)

 
    var response struct {
        AccessToken string `json:"access_token"`
    }
    err := json.NewDecoder(w.Body).Decode(&response)
    assert.NoError(t, err)
 
 
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
	assert.JSONEq(t, `{"error": "Usuario no existe!"}`, w.Body.String())
}

// pruebas Logout 

func TestLogout(t *testing.T) {
  
    tx := initializer.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            tx.Rollback()  
        }
    }()

 
	user := createUser(t, "josenaranjo", "password123")

 
    r := setupRouter()
    reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email, user.Password)
    req, _ := http.NewRequest(http.MethodPost, "/logout", strings.NewReader(reqBody))
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

 

    var revokedToken models.RevokedToken
	if err := tx.Raw("SELECT * FROM revoked_tokens WHERE token = ?", token).Scan(&revokedToken).Error; err != nil {
		t.Errorf("Failed to find revoked token: %v", err)
	}

    if err := tx.Unscoped().Delete(&user).Error; err != nil {
        t.Errorf("Failed to delete user: %v", err)
    }
	if err := tx.Unscoped().Where("token = ?", token).Delete(&revokedToken).Error; err != nil {
		t.Errorf("Failed to delete revoked token: %v", err)
	}
}

//!!
func TestDeleteOldRecords(t *testing.T) {
  
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to create test database: %v", err)
    }

 
    err = db.AutoMigrate(&models.RevokedToken{})
    if err != nil {
        t.Fatalf("failed to migrate RevokedToken model: %v", err)
    }

    
    oldRecord := models.RevokedToken{Token: "old_token"}
    newRecord := models.RevokedToken{Token: "new_token"}

    db.Create(&oldRecord)
    time.Sleep(10 * time.Second)  
    db.Create(&newRecord)
 
    util.DeleteExpiredTokens(db)
 
    var count int64
    db.Model(&models.RevokedToken{}).Where("token = ?", oldRecord.Token).Count(&count)
    if count != 0 {
        t.Errorf("expected old record to be deleted, but it still exists")
    }
 
    db.Model(&models.RevokedToken{}).Where("token = ?", newRecord.Token).Count(&count)
    if count != 1 {
        t.Errorf("expected new record to be kept, but it was deleted")
    }
}



func TestRevokeToken(t *testing.T){
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to create test database: %v", err)
    }

    err = db.AutoMigrate(&models.RefreshToken{})
    if err != nil {
        t.Fatalf("failed to migrate RevokedToken model: %v", err)
    }

    userID := uint(1)

    c,_ := gin.CreateTestContext(nil)

    refreshToken,err := util.GenerateRefreshToken(c,userID)

    if err != nil{
        t.Fatalf("failed to generate refresh token: %v", err)
    }
    if refreshToken == "" {
        t.Fatal("expected refresh token to be not empty")
    }


    existingToken := &models.RefreshToken{
        UserID:    userID,
        Token:     refreshToken,
        ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
    }
    db.Create(existingToken)

    refreshToken, err = util.GenerateRefreshToken(c, userID)
    if err != nil {
        t.Fatalf("failed to generate refresh token: %v", err)
    }
    if refreshToken != existingToken.Token {
        t.Fatal("expected refresh token to match existing token")
    }
}

func TestRenewSession(t*testing.T){
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to create test database: %v", err)
    }

    err = db.AutoMigrate(&models.User{})
    if err != nil {
        t.Fatalf("failed to migrate User model: %v", err)
    }

    err = db.AutoMigrate(&models.RefreshToken{})
    if err != nil {
        t.Fatalf("failed to migrate RefreshToken model: %v", err)
    }

user := createUser(t,"test@test.com","password123")

refreshToken := models.RefreshToken{
    Token:     "test_refresh_token",
    UserID:    uint(user.ID),
    ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
}
db.Create(&refreshToken)
c, _ := gin.CreateTestContext(nil)
var jsonStr = []byte(`{"refresh_token": "test_refresh_token"}`)
c.Request, _ = http.NewRequest(http.MethodPost, "/refresh_token", bytes.NewBuffer(jsonStr))
c.Request.Header.Set("Content-Type", "application/json")


RefreshToken(c)
if c.Writer.Status() != http.StatusOK {
    t.Fatalf("expected status code %d, but got %d", http.StatusOK, c.Writer.Status())
}


var response gin.H
if err := c.ShouldBindJSON(&response); err != nil {
    t.Fatalf("failed to parse response body: %v", err)
}
if _, ok := response["token"]; !ok {
    t.Fatal("expected response body to contain 'token' field")
}
}