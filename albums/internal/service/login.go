package service

import (
	"albums/internal/models"
	"crypto/rand"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginService interface {
	JWTTokenGenerator(userID string) (string, error)
	GetUserbyID(userID string) (models.Users, error)
	GetUserInfo(username string) (models.Users, error)
	RegisterUser(username, password string) (models.Users, error)
	ValidateCredentials(username, password string) (bool, uint, error)
}

type LoginServiceImpl struct {
	DB *gorm.DB
}

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

// func addTokenToDB(userID, token string) error {
// 	return database.DB.Model(models.Users{}).Where("id = ?", userID).Update("token", token).Error
// }

func (service *LoginServiceImpl) JWTTokenGenerator(userID string) (string, error) {

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

func (service *LoginServiceImpl) GetUserbyID(userID string) (models.Users, error) {
	var user models.Users
	if err := service.DB.Find(&user, userID).Error; err != nil {
		return models.Users{}, err
	}
	return user, nil
}

func (service *LoginServiceImpl) GetUserInfo(username string) (models.Users, error) {
	var user models.Users
	if err := service.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (service *LoginServiceImpl) ValidateCredentials(username, password string) (bool, uint, error) {
	user, err := service.GetUserInfo(username)
	if err != nil {
		return false, 0, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, 0, err
	}
	return true, user.ID, nil
}

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

func (service *LoginServiceImpl) RegisterUser(username, password string) (models.Users, error) {
	_, err := service.GetUserInfo(username)
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
	if err := service.DB.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return models.Users{}, errors.New("user with same username exists")
		}
		return models.Users{}, err
	}
	return user, nil
}
