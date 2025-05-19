package service

import (
	"albums/internal/database"
	"albums/internal/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var TokenStore = make(map[string]int64)

func TokenCleaner(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			for token, expiryTime := range TokenStore {
				if time.Now().Unix() > expiryTime {
					delete(TokenStore, token)
				}
			}
		}
	}()
}

func TokenGenerator() string {
	token := uuid.NewString()
	expiryTime := time.Now().Add(time.Minute * 5).Unix()
	TokenStore[token] = expiryTime
	return token
}

func addTokenToDB(userID, token string) error {
	return database.DB.Model(models.Users{}).Where("id = ?", userID).Update("token", token).Error
}

func JWTTokenGenerator(userID string) (string, error) {

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	jwtExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINUTES"))
	expiryTime := time.Now().Add(time.Minute * time.Duration(jwtExpiryMinutes))
	claims := TokenClaim{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	// if err := addTokenToDB(userID, signedToken); err != nil {
	// 	return "", err
	// }
	return signedToken, err
}

// func JWTTokenValidator(bearerToken string) error {
// 	var token *jwt.Token
// 	claims := &tokenClaim{}
// 	token, err := jwt.ParseWithClaims(bearerToken, claims, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("JWT_SECRET")), nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func getUserbyID(userID string) (models.Users, error) {
	var user models.Users
	if err := database.DB.Find(&user, userID).Error; err != nil {
		return models.Users{}, err
	}
	return user, nil
}

func getUserInfo(username string) (models.Users, error) {
	var user models.Users
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func ValidateCredentials(username, password string) (bool, uint, error) {
	user, err := getUserInfo(username)
	if err != nil {
		return false, 0, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, 0, err
	}
	return true, user.ID, nil
}
