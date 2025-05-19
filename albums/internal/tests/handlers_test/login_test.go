package handlers_test

import (
	"albums/internal/handlers"
	"albums/internal/mocks"
	"albums/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    models.Users
		mockResponse   models.Users
		mockError      error
		expectedCode   int
		expectedResult string
	}{
		{
			name:           "successfully register user",
			requestBody:    models.Users{Username: "john", Password: "secret"},
			mockResponse:   models.Users{Username: "john", Password: "hashed"},
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedResult: `user created`,
		},
		{
			name:           "duplicate username",
			requestBody:    models.Users{Username: "john", Password: "secret"},
			mockResponse:   models.Users{},
			mockError:      assert.AnError,
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `assert.AnError general error for testing`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.LoginService)
			mockService.On("RegisterUser", tt.requestBody.Username, tt.requestBody.Password).
				Return(tt.mockResponse, tt.mockError)

			handler := handlers.LoginHandler{LoginService: mockService}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Register(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockToken      string
		mockError      error
		expectedCode   int
		expectedResult string
	}{
		{
			name:           "successful login",
			userID:         "123",
			mockToken:      "token123",
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedResult: `Loggedin`,
		},
		{
			name:           "user not found in context",
			userID:         "",
			expectedCode:   http.StatusNotFound,
			expectedResult: `user not found`,
		},
		{
			name:           "token generation failed",
			userID:         "123",
			mockError:      assert.AnError,
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `Not able to generate token`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.LoginService)

			if tt.userID != "" {
				mockService.On("JWTTokenGenerator", tt.userID).
					Return(tt.mockToken, tt.mockError)
			}

			handler := handlers.LoginHandler{LoginService: mockService}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/login", nil)

			// Simulate middleware injecting UserID
			if tt.userID != "" {
				c.Set("UserID", tt.userID)
			}
			handler.Login(c)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}
