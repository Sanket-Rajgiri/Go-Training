package models

import (
	"gorm.io/gorm"
)

type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type Users struct {
	gorm.Model
	Username  string `json:"username" binding:"required" gorm:"unique;not null;default null"`
	Password  string `json:"password" binding:"required"`
	SecretKey string
	Role      string
	Revoked   bool `gorm:"default:false"`
}

type UserSwagger struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
