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


func TestRefreshToken(t *testing.T) {
	// Setup
	router := gin.Default()
	router.POST("/refresh-token", controllers.RefreshToken)

	// Mock database setup (you'll need to implement this)
	// mockDB := setupMockDB()
	// initializer.DB = mockDB

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Invalid request body",
			requestBody:    map[string]string{},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "Invalid request body"},
		},
		{
			name:        "Invalid refresh token",
			requestBody: map[string]string{"refresh_token": "invalid_token"},
			setupMock: func() {
				// mockDB.ExpectQuery(...).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Invalid refresh token"},
		},
		{
			name:        "Expired refresh token",
			requestBody: map[string]string{"refresh_token": "eyJLZXlJZCI6InNlY3JldCIsImFsZyI6IkhTMjU2IiwidHlwIjoiSldUIn0.eyJhdWQiOiJodHRwOi8vbG9jYWxob3N0OjUwMDAiLCJleHAiOjE0MjIyNzA5MjksImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODA4MCIsInN1YiI6IjE1OCIsInVzZXJpZCI6IjE1OCIsInVzZXJuYW1lIjoibWlsdG9uX0BzdXBlcmFkbWluLmNvbSIsInVzZXJ0eXBlIjoic3VwZXJhZG1pbiJ9.qGxLNwd3yqsNRm11TSDihvrpVE8hUEZiDf0sHRNp6DI"},
			setupMock: func() {
				// mockDB.ExpectQuery(...).WillReturnRows(sqlmock.NewRows([]string{"token", "expires_at", "user_id"}).
				//     AddRow("expired_token", time.Now().Add(-24*time.Hour), 1))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Refresh token expired"},
		},
		{
			name:        "Successful token refresh",
			requestBody: map[string]string{"refresh_token": "77589c02-0e0d-4a44-8c54-99a9b09a5e1a"},
			setupMock: func() {
				// mockDB.ExpectQuery(...).WillReturnRows(sqlmock.NewRows([]string{"token", "expires_at", "user_id"}).
				//     AddRow("valid_token", time.Now().Add(24*time.Hour), 1))
				// Add expectations for user query, token generation, etc.
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"session_token": "new_session_token", "user_token": "new_refresh_token"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.setupMock()

			// Create request
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/refresh-token", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Check response body
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			for key, value := range tc.expectedBody {
				assert.Equal(t, value, response[key])
			}

			// Assert that all expectations were met
			// assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}