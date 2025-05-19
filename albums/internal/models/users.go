package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Username  string `json:"username" binding:"required" gorm:"unique;not null;default null"`
	Password  string `json:"password" binding:"required"`
	SecretKey string
	Role      string
}
