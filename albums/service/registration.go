package service

import (
	"albums/database"
	"albums/models"
	"crypto/rand"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type TokenClaim struct {
	Role string `json:"Role"`
	jwt.RegisteredClaims
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = charset[int(b)%len(charset)]
	}

	return string(bytes), nil
}

func RegisterUser(username, password string) (models.Users, error) {
	_, err := getUserInfo(username)
	if err == nil {
		return models.Users{}, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Users{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Users{}, err
	}

	secretKey, err := generateRandomString(16)
	if err != nil {
		return models.Users{}, err
	}

	user := models.Users{Username: username, Password: string(hashedPassword), SecretKey: secretKey, Role: "user"}
	if err := database.DB.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return models.Users{}, errors.New("user with same username exists")
		}
		return models.Users{}, err
	}
	return user, nil
}
