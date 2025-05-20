package handlers_test

import (
	"albums/internal/handlers"
	"albums/internal/mocks"
	"albums/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAlbums(t *testing.T) {
	type response struct {
		Albums []models.Album `json:"albums"`
	}

	tests := []struct {
		name         string
		mockResponse []models.Album
		mockError    error
		expectedCode int
	}{
		{
			name: "Success",
			mockResponse: []models.Album{
				{
					Model:  gorm.Model{ID: 1}, // if you're using gorm.Model
					Title:  "Test Album",
					Artist: "Test Artist",
					Price:  99.99,
				},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "InternalError",
			mockResponse: nil,
			mockError:    errors.New("db error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockService := new(mocks.AlbumService)
			mockService.On("GetAlbums").Return(tt.mockResponse, tt.mockError)

			handler := handlers.AlbumHandler{AlbumService: mockService}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/albums", nil)

			handler.GetAlbums(c)
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.mockError == nil {
				var resp response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp.Albums, len(tt.mockResponse))

				for i, album := range tt.mockResponse {
					assert.Equal(t, album.ID, resp.Albums[i].ID)
					assert.Equal(t, album.Title, resp.Albums[i].Title)
					assert.Equal(t, album.Artist, resp.Albums[i].Artist)
					assert.Equal(t, album.Price, resp.Albums[i].Price)
				}
			} else {
				assert.Contains(t, w.Body.String(), `"error": "`+tt.mockError.Error()+`"`)
			}
		})
	}
}

func TestGetAlbumByID(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		mockResponse models.Album
		mockError    error
		expectedCode int
	}{
		{
			name:    "Success",
			paramID: "1",
			mockResponse: models.Album{
				Model:  gorm.Model{ID: 1}, // if you're using gorm.Model
				Title:  "Test Album",
				Artist: "Test Artist",
				Price:  99.99,
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid ID",
			paramID:      "abc",
			mockError:    errors.New("Invalid ID"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "InternalError",
			paramID:      "2",
			mockResponse: models.Album{},
			mockError:    errors.New("db error"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockService := new(mocks.AlbumService)
			if id, err := strconv.Atoi(tt.paramID); err == nil {
				mockService.On("GetAlbumByID", uint(id)).Return(tt.mockResponse, tt.mockError)
			}

			handler := handlers.AlbumHandler{AlbumService: mockService}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/albums/"+tt.paramID, nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.paramID}}

			handler.GetAlbumByID(c)
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.mockError == nil {
				var resp models.Album
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse.ID, resp.ID)
				assert.Equal(t, tt.mockResponse.Title, resp.Title)
				assert.Equal(t, tt.mockResponse.Artist, resp.Artist)
				assert.Equal(t, tt.mockResponse.Price, resp.Price)
			} else {
				assert.Contains(t, w.Body.String(), `"error": "`+tt.mockError.Error()+`"`)
			}
		})
	}
}

func TestAddAlbums(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		mockInput      models.Album
		mockReturnID   uint
		mockError      error
		expectedCode   int
		expectedResult string
	}{
		{
			name:           "Success",
			requestBody:    `{"title":"New","artist":"Artist","price":15.5}`,
			mockInput:      models.Album{Title: "New", Artist: "Artist", Price: 15.5},
			mockReturnID:   101,
			mockError:      nil,
			expectedCode:   http.StatusCreated,
			expectedResult: `"message": "Album Added"`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{invalid json}`,
			expectedCode:   http.StatusBadRequest,
			expectedResult: `"error": "Invalid JSON"`,
		},
		{
			name:           "Service Error",
			requestBody:    `{"title":"Bad","artist":"X","price":1}`,
			mockInput:      models.Album{Title: "Bad", Artist: "X", Price: 1},
			mockReturnID:   0,
			mockError:      errors.New("db fail"),
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `"error": "db fail"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.AlbumService)
			if tt.mockError != nil || tt.mockReturnID != 0 {
				mockService.On("AddAlbums", tt.mockInput).Return(tt.mockReturnID, tt.mockError)
			}

			handler := handlers.AlbumHandler{AlbumService: mockService}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodPost, "/albums", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.AddAlbums(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}

func TestUpdatePrice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		mockInput      models.Album
		mockError      error
		expectedCode   int
		expectedResult string
	}{
		{
			name:        "Success",
			requestBody: `{"id":1,"price":20.5}`,
			mockInput: models.Album{
				Model: gorm.Model{ID: 1}, Price: 20.5},
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedResult: `"message": "Album Price Updated Successfully"`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `bad json`,
			expectedCode:   http.StatusBadRequest,
			expectedResult: `"error": "Invalid JSON"`,
		},
		{
			name:        "Service Error",
			requestBody: `{"id":2,"price":11}`,
			mockInput: models.Album{
				Model: gorm.Model{ID: 2}, Price: 11},
			mockError:      errors.New("update fail"),
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `"error": "update fail"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.AlbumService)
			if tt.mockError != nil || tt.mockInput.ID != 0 {
				mockService.On("UpdatePrice", tt.mockInput).Return(tt.mockInput.ID, tt.mockError)
			}

			handler := handlers.AlbumHandler{AlbumService: mockService}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodPut, "/albums", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.UpdatePrice(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}

func TestDeleteAlbum(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		paramID        string
		mockReturnID   uint
		mockError      error
		expectedCode   int
		expectedResult string
	}{
		{
			name:           "Success",
			paramID:        "1",
			mockReturnID:   1,
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedResult: `"message": "album deleted successfully"`,
		},
		{
			name:           "Invalid ID",
			paramID:        "abc",
			expectedCode:   http.StatusBadRequest,
			expectedResult: `"error": "Invalid ID"`,
		},
		{
			name:           "Service Error",
			paramID:        "2",
			mockReturnID:   0,
			mockError:      errors.New("not found"),
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `"error": "not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.AlbumService)
			if tt.mockError != nil || tt.mockReturnID != 0 {
				id, _ := strconv.Atoi(tt.paramID)
				mockService.On("DeleteAlbum", uint(id)).Return(tt.mockReturnID, tt.mockError)
			}

			handler := handlers.AlbumHandler{AlbumService: mockService}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.paramID}}
			req := httptest.NewRequest(http.MethodDelete, "/albums/"+tt.paramID, nil)
			c.Request = req

			handler.DeleteAlbum(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedResult)
		})
	}
}
