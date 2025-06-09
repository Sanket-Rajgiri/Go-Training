package routes_test

import (
	"albums/internal/config/env"
	"albums/internal/handlers"
	"albums/internal/models"
	"albums/internal/service"
	testutils "albums/internal/testUtils"
	"albums/routes"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func initLoginSetup(db *gorm.DB) *gin.Engine {
	loginService := &service.LoginServiceImpl{DB: db}

	loginHandler := &handlers.LoginHandler{
		LoginService: loginService,
	}

	router := gin.New()
	router.RedirectTrailingSlash = false
	routes.RegisterLoginRoutes(router, loginHandler, loginService)
	return router
}

type loginTest struct {
	name             string
	body             interface{}
	setup            func(db *gorm.DB)
	expectedStatus   int
	expectedContains string
	token            string
	requestPath      string
}

func TestRegisterIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initLoginSetup(db)

	tests := []loginTest{
		{
			name: "valid request",
			body: models.Users{
				Username: "sanket",
				Password: "password",
			},
			requestPath:      "/login/register",
			expectedStatus:   http.StatusOK,
			expectedContains: "user created",
		},
		{
			name:             "Invalid JSON",
			body:             `{"body":"invalid `,
			requestPath:      "/login/register",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid JSON",
		},
		{
			name: "Username Already Exists",
			body: models.Users{
				Username: "test",
				Password: "testpass",
			},
			setup: func(db *gorm.DB) {
				db.Create(&models.Users{
					Username: "test",
					Password: "different",
				})
			},
			requestPath:      "/login/register",
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "username already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}
			var req *http.Request
			switch v := tt.body.(type) {
			case string:
				req, err = http.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBuffer(jsonBody))
			}
			req.Header.Set("Authorization", "Bearer "+tt.token)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedContains != "" {
				assert.Contains(t, w.Body.String(), tt.expectedContains)
			}
		})
	}
}

func TestLogoutIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initLoginSetup(db)
	env.JWT_SECRET = "test"
	tests := []loginTest{
		{
			name: "valid request",
			body: `{"username":"test-user"}`,
			setup: func(db *gorm.DB) {
				db.Create(&models.Users{
					Username:  "test-user",
					Password:  "test-pass",
					SecretKey: "randomstring",
					Role:      "reader",
				})
			},
			expectedStatus:   http.StatusOK,
			expectedContains: "Logged out",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
			requestPath:      "/logout",
		},
		{
			name:             "username not exist",
			body:             `{"username":"test-user-2"}`,
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "record not found",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
			requestPath:      "/logout",
		},
		{
			name:             "Invalid Body",
			body:             `{"username":"test-user-3"`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid JSON",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
			requestPath:      "/logout",
		},
		{
			name:             "Invalid Token",
			body:             `{"username":"test-user-3"}`,
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "token has invalid claims",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkzODIyNzQsImV4cCI6MTc0OTM4MjQwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSJ9.lZ-fGjFVIOyOwv3vb6qtueEbTdvsnRh3ygkKFtafstc",
			requestPath:      "/logout",
		},
		{
			name:             "Missing Token",
			body:             `{"username":"test-user-3"}`,
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Empty token",
			requestPath:      "/logout",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}
			var req *http.Request
			switch v := tt.body.(type) {
			case string:
				req, err = http.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBuffer(jsonBody))
			}
			req.Header.Set("Authorization", "Bearer "+tt.token)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedContains != "" {
				assert.Contains(t, w.Body.String(), tt.expectedContains)
			}
		})
	}
}

func TestLoginIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initLoginSetup(db)
	env.JWT_SECRET = "test"
	tests := []loginTest{
		{
			name: "valid request",
			body: `{"username": "test-user","password":"test-pass"}`,
			setup: func(db *gorm.DB) {
				pass := "test-pass"
				hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
				db.Create(&models.Users{
					Username:  "test-user",
					Password:  string(hashedPass),
					SecretKey: "randomString",
					Role:      "reader",
				})
			},
			requestPath:      "/login/",
			expectedStatus:   http.StatusOK,
			expectedContains: "Loggedin",
		},
		{
			name:             "Invalid Body",
			body:             `{"username": "test-user2","password":"test-pass"`,
			requestPath:      "/login/",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "invalid json",
		},
		{
			name:             "username not found",
			body:             `{"username": "test-user2","password":"test-pass2"}`,
			requestPath:      "/login/",
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "username not found",
		},
		{
			name: "Invalid Password",
			body: `{"username": "test-user2","password":"test-pass2"}`,
			setup: func(db *gorm.DB) {
				pass := "test-pass"
				hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
				db.Create(&models.Users{
					Username:  "test-user2",
					Password:  string(hashedPass),
					SecretKey: "randomString",
					Role:      "reader",
				})
			},
			requestPath:      "/login/",
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "invalid credentials",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}
			var req *http.Request
			switch v := tt.body.(type) {
			case string:
				req, err = http.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBuffer(jsonBody))
			}
			req.Header.Set("Authorization", "Bearer "+tt.token)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedContains != "" {
				assert.Contains(t, w.Body.String(), tt.expectedContains)
			}
		})
	}
}

func TestRefreshTokenIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initLoginSetup(db)
	env.JWT_SECRET = "test"
	tests := []loginTest{
		{
			name:        "valid request",
			requestPath: "/login/refresh",
			body: map[string]string{
				"username": "test-user",
				"token":    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkzODIyNzQsImV4cCI6MTc4MDk4OTI3NywiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMiIsIlJvbGUiOiJyZWFkZXIifQ.wfHaA4TIyG_X0AG8E-2VAnRFoMtSTa6LHdRSzb8vKOE",
			},
			setup: func(db *gorm.DB) {
				pass := "test-pass"
				hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
				db.Create(&models.Users{
					Model:     gorm.Model{ID: 2},
					Username:  "test-user",
					Password:  string(hashedPass),
					SecretKey: "randomString",
					Role:      "reader",
				})
			},
			expectedStatus:   http.StatusOK,
			expectedContains: "refreshed",
		},
		{
			name:             "Invalid Body",
			requestPath:      "/login/refresh",
			body:             `{"username":"admin-app", "password":"admin@123"}`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid JSON",
		},
		{
			name:        "username not found",
			requestPath: "/login/refresh",
			body: map[string]string{
				"username": "test-user2",
				"token":    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkzODIyNzQsImV4cCI6MTc4MDk4OTI3NywiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMiIsIlJvbGUiOiJyZWFkZXIifQ.wfHaA4TIyG_X0AG8E-2VAnRFoMtSTa6LHdRSzb8vKOE",
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "user not found",
		},
		{
			name:        "Invalid Token Signature",
			requestPath: "/login/refresh",
			body: map[string]string{
				"username": "test-user",
				"token":    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkzODIyNzQsImV4cCI6MTc4MDk4OTI3NywiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMiIsIlJvbGUiOiJyZWFkZXIifQ.VuWm4IRRqpbItIQa-aEJ7Ib4GzZfmvmk_v57r1LNaAg",
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "token signature is invalid",
		},
		{
			name:        "Invalid Token Claim",
			requestPath: "/login/refresh",
			body: map[string]string{
				"username": "test-user",
				"token":    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkzODIyNzQsImV4cCI6MTc4MDk4OTI3NywiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMiIsIlJvbGUiOiJhZG1pbiJ9.9Ja_1TWxdzVKvOR1-VmAg8pML9YaXyEOoHrmVGLJLBs",
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "invalid token claim",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}
			var req *http.Request
			switch v := tt.body.(type) {
			case string:
				req, err = http.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBuffer(jsonBody))
			}
			req.Header.Set("Authorization", "Bearer "+tt.token)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedContains != "" {
				assert.Contains(t, w.Body.String(), tt.expectedContains)
			}
		})
	}
}
