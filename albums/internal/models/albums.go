package models

import (
	"time"

	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Title  string  `json:"title" binding:"required" example:"Blue Train"`
	Artist string  `json:"artist" binding:"required" example:"John Coltrane"`
	Price  float64 `json:"price" binding:"required" example:"56.99"`
}

// AlbumSwagger is used for Swagger documentation only
type AlbumSwagger struct {
	ID        uint       `json:"id" example:"1"`
	CreatedAt time.Time  `json:"created_at" example:"2023-05-19T15:04:05Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2023-05-19T15:04:05Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" example:"2023-05-20T15:04:05Z"`
	Title     string     `json:"title" example:"Blue Train"`
	Artist    string     `json:"artist" example:"John Coltrane"`
	Price     float64    `json:"price" example:"56.99"`
}
