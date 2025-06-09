package integration_test

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
	"gorm.io/gorm"
)

func initAlbumsSetup(db *gorm.DB) *gin.Engine {
	albumService := &service.AlbumServiceImpl{DB: db}

	albumHandler := &handlers.AlbumHandler{
		AlbumService: albumService,
	}

	router := gin.New()
	router.RedirectTrailingSlash = false
	routes.RegisterAlbumRoutes(router, albumHandler)
	return router
}

type albumTest struct {
	name             string
	body             interface{}
	setup            func(db *gorm.DB)
	expectedStatus   int
	expectedContains string
	token            string
	requestPath      string
}

func TestGetAlbumsIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env.JWT_SECRET = "test"
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initAlbumsSetup(db)
	tests := []albumTest{
		{
			name:           "valid request with empty albums",
			expectedStatus: http.StatusOK,
			token:          "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
		},
		{
			name: "valid request with Albums",
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Title:  "Test Title",
					Artist: "Test Artist",
					Price:  1.0,
				})
			},
			expectedStatus:   http.StatusOK,
			expectedContains: "Test Artist",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
		},
		{
			name:             "missing token",
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Empty token",
		},
		{
			name:             "Invalid Token",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkxOTMzMDksImV4cCI6MTc4MDcyOTMwOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMSIsIlJvbGUiOiJQcm9qZWN0IEFkbWluaXN0cmF0b3IifQ.1k3f9r8rS_op8BclNpjUysdoNiQvSDqY7QQgg1BCwQg",
			expectedStatus:   http.StatusForbidden,
			expectedContains: "Access Denied",
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
				req, err = http.NewRequest(http.MethodGet, "/albums/", bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodGet, "/albums/", bytes.NewBuffer(jsonBody))
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

func TestGetAlbumbyIDIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env.JWT_SECRET = "test"
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initAlbumsSetup(db)
	tests := []albumTest{
		{

			name: "valid request",
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Model:  gorm.Model{ID: 1},
					Title:  "Test Title",
					Artist: "Test Artist",
					Price:  1.0,
				})
			},
			expectedStatus:   http.StatusOK,
			expectedContains: "Test Artist",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
			requestPath:      "/albums/1",
		},
		{
			name:             "missing token",
			requestPath:      "/albums/2",
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Empty token",
		},
		{
			name:             "Invalid Token",
			requestPath:      "/albums/3",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkxOTMzMDksImV4cCI6MTc4MDcyOTMwOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMSIsIlJvbGUiOiJQcm9qZWN0IEFkbWluaXN0cmF0b3IifQ.1k3f9r8rS_op8BclNpjUysdoNiQvSDqY7QQgg1BCwQg",
			expectedStatus:   http.StatusForbidden,
			expectedContains: "Access Denied",
		},
		{
			name:             "Record Not Found",
			requestPath:      "/albums/3",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "record not found",
		},
		{
			name:             "Invalid ID",
			requestPath:      "/albums/abc",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid ID",
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
				req, err = http.NewRequest(http.MethodGet, tt.requestPath, bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodGet, tt.requestPath, bytes.NewBuffer(jsonBody))
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

func TestAddAlbumIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env.JWT_SECRET = "test"
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initAlbumsSetup(db)

	tests := []albumTest{
		{
			name: "valid request",
			body: models.Album{
				Title:  "Integration Test",
				Artist: "Navdeep",
				Price:  10.5,
			},
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1NDQyOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.sc_RTpEB41GVWiMMNOCZ-7-_kwgof9BzS6G3kve1j1g",
			expectedStatus:   http.StatusCreated,
			expectedContains: "Album Added",
		},
		{
			name: "invalid token",
			body: models.Album{
				Title:  "Test 2",
				Artist: "User 2",
				Price:  89.2,
			},
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc0OTEyMTg3MSwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.i0l1Ph9e3Vku5-NVQwv3Vj0OLaeURgh-0b_exG0jtzk",
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "token is expired",
		},
		{
			name: "missing token",
			body: models.Album{
				Title:  "Test 2",
				Artist: "User 2",
				Price:  89.2,
			},
			token:            "",
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Empty token",
		},
		{
			name: "invalid role",
			body: models.Album{
				Title:  "Test 2",
				Artist: "User 2",
				Price:  89.2,
			},
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1ODEwOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Z22L6rUc9dZVIH7BvuAa9EY89vSb2g6iHJyQ_vdc1v4",
			expectedStatus:   http.StatusForbidden,
			expectedContains: "route not allowed to access",
		},
		{
			name:             "invalid JSON",
			body:             `{"title": "Oops",`,
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1NDQyOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.sc_RTpEB41GVWiMMNOCZ-7-_kwgof9BzS6G3kve1j1g",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid JSON",
		},
		{
			name: "missing title",
			body: map[string]interface{}{
				"artist": "NoTitle",
				"price":  9.0,
			},
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1NDQyOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.sc_RTpEB41GVWiMMNOCZ-7-_kwgof9BzS6G3kve1j1g",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid JSON",
		},
		{
			name: "non-numeric price",
			body: map[string]interface{}{
				"title":  "BadPrice",
				"artist": "X",
				"price":  "free",
			},
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1NDQyOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.sc_RTpEB41GVWiMMNOCZ-7-_kwgof9BzS6G3kve1j1g",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var req *http.Request

			switch v := tt.body.(type) {
			case string:
				req, err = http.NewRequest(http.MethodPost, "/albums/", bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodPost, "/albums/", bytes.NewBuffer(jsonBody))
			}
			req.Header.Set("Authorization", "Bearer "+tt.token)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedContains)
		})
	}
}

func TestUpdateAlbumIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initAlbumsSetup(db)

	tests := []albumTest{
		{
			name: "valid request",
			body: models.Album{
				Model: gorm.Model{
					ID: 1,
				},
				Title:  "Integration Test",
				Artist: "Navdeep",
				Price:  10.5,
			},
			token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc4MDY1NDQyOCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.sc_RTpEB41GVWiMMNOCZ-7-_kwgof9BzS6G3kve1j1g",
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Title:  "Integration Test",
					Artist: "Navdeep",
					Price:  9.0,
				})
			},
			expectedStatus:   http.StatusOK,
			expectedContains: "Album Price Updated Successfully",
		},
		{
			name: "invalid token",
			body: models.Album{
				Model: gorm.Model{
					ID: 2,
				},
				Title:  "Test 2",
				Artist: "User 2",
				Price:  89.2,
			},
			token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOjE3NDkxMTU0NjMsImV4cCI6MTc0OTEyMTg3MSwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.i0l1Ph9e3Vku5-NVQwv3Vj0OLaeURgh-0b_exG0jtzk",
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Title:  "Test 2",
					Artist: "User 2",
					Price:  87.8,
				})
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "token is expired",
		},
		{
			name: "missing token",
			body: models.Album{
				Model: gorm.Model{
					ID: 3,
				},
				Title:  "Test 2",
				Artist: "User 2",
				Price:  89.2,
			},
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Title:  "Test 3",
					Artist: "User 3",
					Price:  80,
				})
			},
			token:            "",
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Empty token",
		},
		{
			name: "invalid role",
			body: models.Album{
				Model:  gorm.Model{ID: 4},
				Title:  "Test 4",
				Artist: "User 4",
				Price:  9.2,
			},
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Title:  "Test 4",
					Artist: "User 4",
					Price:  8.0,
				})
			},
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkxOTMzMDksImV4cCI6MTc4MDcyOTMwOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMSIsIlJvbGUiOiJyZWFkZXIifQ.Y6LWhUBeSkkEw21jwPIi5CwBGipQwuMhWThCMRdjpS4",
			expectedStatus:   http.StatusForbidden,
			expectedContains: "route not allowed to access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var req *http.Request
			if tt.setup != nil {
				tt.setup(db)
			}
			switch v := tt.body.(type) {
			case string:
				req, err = http.NewRequest(http.MethodPatch, "/albums/", bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodPatch, "/albums/", bytes.NewBuffer(jsonBody))
			}
			req.Header.Set("Authorization", "Bearer "+tt.token)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedContains)
		})
	}
}

func TestDeleteAlbumIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env.JWT_SECRET = "test"
	ctx := context.Background()
	db, cleanup, err := testutils.InitTestDB(ctx)
	require.NoError(t, err)
	defer cleanup()
	testutils.InitOtel()
	router := initAlbumsSetup(db)
	tests := []albumTest{
		{
			name: "valid request",
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Model:  gorm.Model{ID: 1},
					Title:  "Test Title",
					Artist: "Test Artist",
					Price:  1.0,
				})
			},
			requestPath:      "/albums/1",
			expectedStatus:   http.StatusOK,
			expectedContains: "album deleted successfully",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkxOTMzMDksImV4cCI6MTc4MDcyOTMwOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMSIsIlJvbGUiOiJhZG1pbiJ9.KFTF8JbJkNWEUOYQbGibjHWaq8n0knhzoS9zzxOq1-Q",
		},
		{
			name: "missing token",
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Model:  gorm.Model{ID: 2},
					Title:  "Test Title 2",
					Artist: "Test Artist 2",
					Price:  10.0,
				})
			},
			requestPath:      "/albums/2",
			expectedStatus:   http.StatusUnauthorized,
			expectedContains: "Empty token",
		},
		{
			name:             "Invalid Token",
			requestPath:      "/albums/2",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkxOTMzMDksImV4cCI6MTc4MDcyOTMwOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMSIsIlJvbGUiOiJQcm9qZWN0IEFkbWluaXN0cmF0b3IifQ.1k3f9r8rS_op8BclNpjUysdoNiQvSDqY7QQgg1BCwQg",
			expectedStatus:   http.StatusForbidden,
			expectedContains: "Access Denied",
		},
		{
			name:        "Invalid Role",
			requestPath: "/albums/2",
			token:       "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkxOTMzMDksImV4cCI6MTc4MDcyOTMwOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.LsrFX7rP3sMgGyLLCx0cTxB40I2N8EOHaXfUL1eCfu8",
			setup: func(db *gorm.DB) {
				db.Create(&models.Album{
					Model:  gorm.Model{ID: 2},
					Title:  "Test Title 2",
					Artist: "Test Artist 2",
					Price:  10.0,
				})
			},
			expectedStatus:   http.StatusForbidden,
			expectedContains: "route not allowed to access",
		},
		{
			name:             "Record Not Found",
			requestPath:      "/albums/3",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkxOTMzMDksImV4cCI6MTc4MDcyOTMwOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMSIsIlJvbGUiOiJhZG1pbiJ9.KFTF8JbJkNWEUOYQbGibjHWaq8n0knhzoS9zzxOq1-Q",
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: "record not found",
		},
		{
			name:             "Invalid ID",
			requestPath:      "/albums/abc",
			token:            "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkxOTMzMDksImV4cCI6MTc4MDcyOTMwOSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiMSIsIlJvbGUiOiJhZG1pbiJ9.KFTF8JbJkNWEUOYQbGibjHWaq8n0knhzoS9zzxOq1-Q",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid ID",
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
				req, err = http.NewRequest(http.MethodDelete, tt.requestPath, bytes.NewBufferString(v))
			default:
				jsonBody, _ := json.Marshal(v)
				req, err = http.NewRequest(http.MethodDelete, tt.requestPath, bytes.NewBuffer(jsonBody))
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
