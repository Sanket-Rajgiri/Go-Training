package service

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func JWTTokenGenerator(username string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	jwtExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINUTES"))
	expiryTime := time.Now().Add(time.Minute * time.Duration(jwtExpiryMinutes))
	RegisteredClaims := jwt.RegisteredClaims{
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiryTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, RegisteredClaims)
	return token.SignedString(jwtSecret)
}

func JWTTokenValidator(bearerToken string) error {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return errors.New("token expired")
	}
	return nil
}
