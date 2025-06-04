package handlers

import (
	"albums/internal/mocks"
	"albums/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/login/register", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			mockService := new(mocks.LoginService)
			mockService.On("RegisterUser", mock.Anything, tt.requestBody.Username, tt.requestBody.Password).
				Return(tt.mockResponse, tt.mockError)

			handler := LoginHandler{LoginService: mockService}
			handler.Register(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		role             string
		mockAccessToken  string
		mockRefreshToken string
		mockError        error
		expectedCode     int
		expectedResult   string
	}{
		{
			name:             "successful login",
			userID:           "123",
			role:             "user",
			mockAccessToken:  "token123",
			mockRefreshToken: "token90976",
			mockError:        nil,
			expectedCode:     http.StatusOK,
			expectedResult:   `Loggedin`,
		},
		{
			name:           "user not found in context",
			userID:         "",
			expectedCode:   http.StatusNotFound,
			expectedResult: `user not found`,
		},
		{
			name:           "role not found in context",
			userID:         "123",
			role:           "",
			expectedCode:   http.StatusNotFound,
			expectedResult: `role not found`,
		},
		{
			name:           "token generation failed",
			userID:         "123",
			role:           "user",
			mockError:      assert.AnError,
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `assert.AnError general error for testing`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/login", nil)

			// Simulate middleware injecting UserID
			if tt.userID != "" {
				c.Set("UserID", tt.userID)
			}
			if tt.role != "" {
				c.Set("Role", tt.role)
			}

			mockService := new(mocks.LoginService)

			if tt.userID != "" && tt.role != "" {
				mockService.On("Login", mock.Anything, tt.userID, tt.role).
					Return(tt.mockAccessToken, tt.mockRefreshToken, tt.mockError)
			}

			handler := LoginHandler{LoginService: mockService}
			handler.Login(c)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      RefreshHandlerBody
		mockRefreshToken string
		mockError        error
		expectedCode     int
		expectedResult   string
	}{
		{
			name: "successful refresh",
			requestBody: RefreshHandlerBody{
				Username: "abc",
				Token:    "swdwcax",
			},
			mockRefreshToken: "token343e3242",
			mockError:        nil,
			expectedCode:     http.StatusOK,
			expectedResult:   `refreshed`,
		},
		{
			name: "Invalid Body",
			requestBody: RefreshHandlerBody{
				Username: "abc",
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: `Invalid JSON`,
		},
		{
			name: "token generation failed",
			requestBody: RefreshHandlerBody{
				Username: "xyz",
				Token:    "w234ty54rdw",
			},
			mockError:      assert.AnError,
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `assert.AnError general error for testing`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			jsonBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Errorf("error in parsing struct : %v", err)
			}
			c.Request = httptest.NewRequest(http.MethodPost, "/login/refresh", strings.NewReader(string(jsonBody)))
			c.Request.Header.Set("Content-Type", "application/json")

			mockService := new(mocks.LoginService)
			mockService.On("RefreshToken", mock.Anything, tt.requestBody.Username, tt.requestBody.Token).Return(tt.mockRefreshToken, tt.mockError)
			handler := LoginHandler{LoginService: mockService}
			handler.Refresh(c)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockUsername   string
		mockError      error
		expectedCode   int
		expectedResult string
	}{
		{
			name: "successful Logout",
			requestBody: map[string]interface{}{
				"username": "abc",
			},
			mockUsername:   "abc",
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedResult: `Logged out`,
		},
		{
			name: "token generation failed",
			requestBody: map[string]interface{}{
				"username": "123",
			},
			mockUsername:   "123",
			mockError:      assert.AnError,
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `assert.AnError general error for testing`,
		},
		{
			name:           "Invalid Body",
			requestBody:    nil,
			expectedCode:   http.StatusBadRequest,
			expectedResult: `'username' is required`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			jsonBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Errorf("error in parsing struct : %v", err)
			}
			c.Request = httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(string(jsonBody)))
			c.Request.Header.Set("Content-Type", "application/json")

			mockService := new(mocks.LoginService)
			if len(tt.mockUsername) > 0 {
				mockService.On("Logout", mock.Anything, tt.mockUsername).Return(tt.mockError)
			}
			handler := LoginHandler{LoginService: mockService}
			handler.Logout(c)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}
