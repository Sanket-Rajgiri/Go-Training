package middleware_test

import (
	"albums/internal/config/env"
	"albums/internal/customlogs"
	"albums/internal/middleware"
	"albums/internal/mocks"
	"albums/internal/models"
	"albums/internal/service"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	noopTrace "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

func initOtelSetup() {
	otel.SetTracerProvider(noopTrace.NewTracerProvider())
	customlogs.OtelLogger = otelzap.New(zap.NewNop())
}

func TestRoleValidationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		role           interface{}
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid role",
			role:           "user",
			method:         http.MethodPost,
			path:           "/albums/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Ok",
		},
		{
			name:           "role not found",
			role:           nil,
			method:         http.MethodPost,
			path:           "/login",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "role not found",
		},
		{
			name:           "invalid role type",
			role:           make(map[string]string),
			method:         http.MethodPost,
			path:           "/albums/",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "invalid role type",
		},
		{
			name:           "unknown role",
			role:           "guest",
			method:         http.MethodPost,
			path:           "/albums/",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Access Denied",
		},
		{
			name:           "route not allowed",
			role:           "user",
			method:         http.MethodGet,
			path:           "/albums/",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "route not allowed to access",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initOtelSetup()
			router := gin.New()
			router.Use(func(c *gin.Context) {
				if tt.role != nil {
					c.Set("Role", tt.role)
				}
				c.Next()
			})
			router.Use(middleware.RoleValidationMiddleware())

			router.POST("/albums/", func(c *gin.Context) {
				c.String(http.StatusOK, "Ok")
			})

			router.GET("/invalid", func(c *gin.Context) {
				c.String(http.StatusOK, "invalid")
			})

			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}

}

func TestDBAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		requestBody    models.LoginPayload
		mockUserID     uint
		mockRole       string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful auth",
			requestBody: models.LoginPayload{
				Username: "validUsername",
				Password: "validPassword",
			},
			mockUserID:     1,
			mockRole:       "user",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Ok",
		},
		{
			name: "unsuccessful auth",
			requestBody: models.LoginPayload{
				Username: "validUsername",
				Password: "validPassword",
			},
			mockUserID:     0,
			mockRole:       " ",
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "assert.AnError general error for testing",
		},
		{
			name: "empty username and password",
			requestBody: models.LoginPayload{
				Username: "",
				Password: "",
			},
			mockUserID:     0,
			mockRole:       "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid json",
		},
		{
			name: "internal error during validation",
			requestBody: models.LoginPayload{
				Username: "errorUser",
				Password: "password",
			},
			mockUserID:     0,
			mockRole:       "",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error",
		},
		{
			name: "invalid credentials",
			requestBody: models.LoginPayload{
				Username: "invalidUsername",
				Password: "wrongPassword",
			},
			mockUserID:     0,
			mockRole:       "",
			mockError:      errors.New("invalid credentials"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "invalid credentials",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initOtelSetup()

			mockService := new(mocks.LoginService)
			mockService.On("ValidateCredentials", mock.Anything, tt.requestBody.Username, tt.requestBody.Password).
				Return(tt.mockUserID, tt.mockRole, tt.mockError)

			router := gin.New()
			router.Use(middleware.DBAuthMiddleware(mockService))
			router.GET("/login", func(c *gin.Context) {
				c.String(http.StatusOK, "Ok")
			})

			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodGet, "/login", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func generateTestJWT(role string) string {
	claims := &service.TokenClaim{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "1",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(env.JWT_SECRET))
	return signedToken
}

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validToken := generateTestJWT("admin")
	router := gin.New()
	router.Use(middleware.JwtAuthMiddleware())

	// A protected route
	router.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing or invalid Authorization header",
		},
		{
			name:           "invalid auth scheme",
			authHeader:     "Token some-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing or invalid Authorization header",
		},
		{
			name:           "empty bearer token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Empty token",
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid.token.value",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "token is malformed",
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initOtelSetup()
			req, _ := http.NewRequest(http.MethodGet, "/login", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestBasicAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup routes
	router := gin.New()
	router.Use(middleware.BasicAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "authenticated")
	})

	tests := []struct {
		name           string
		username       string
		password       string
		setAuthHeader  bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing auth header",
			setAuthHeader:  false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized",
		},
		{
			name:           "invalid credentials",
			username:       "wronguser",
			password:       "wrongpass",
			setAuthHeader:  true,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized",
		},
		{
			name:           "valid credentials",
			username:       "user",  // make sure this exists in your BasicAuthAccounts map
			password:       "admin", // and the password matches
			setAuthHeader:  true,
			expectedStatus: http.StatusOK,
			expectedBody:   "authenticated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.setAuthHeader {
				auth := tt.username + ":" + tt.password
				basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
				req.Header.Set("Authorization", basicAuth)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
