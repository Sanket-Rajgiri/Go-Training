package service

import (
	"albums/internal/customlogs"
	"albums/internal/models"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginService interface {
	JWTTokenGenerator(ctx context.Context, userID string) (string, error)
	GetUserbyID(userID string) (models.Users, error)
	GetUserInfo(ctx context.Context, username string) (models.Users, error)
	RegisterUser(ctx context.Context, username, password string) (models.Users, error)
	ValidateCredentials(ctx context.Context, username, password string) (bool, uint, error)
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

func (service *LoginServiceImpl) JWTTokenGenerator(ctx context.Context, userID string) (string, error) {
	tracer := otel.Tracer("Service-Tracer")
	_, span := tracer.Start(ctx, "JWTTokenGenerator")
	defer span.End()
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
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	// if err := addTokenToDB(userID, signedToken); err != nil {
	// 	return "", err
	// }
	customlogs.OtelLogger.Info(fmt.Sprintf("userID: %s token generated Successfully", userID))
	span.SetStatus(codes.Ok, "token generated")
	span.SetAttributes(attribute.String("userID", userID))
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

func (service *LoginServiceImpl) GetUserInfo(ctx context.Context, username string) (models.Users, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "GetUserInfo")
	defer span.End()
	var user models.Users
	if err := service.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return user, err
	}
	span.SetStatus(codes.Ok, "user fetched")
	span.SetAttributes(attribute.Int64("userID", int64(user.ID)))
	return user, nil
}

func (service *LoginServiceImpl) ValidateCredentials(ctx context.Context, username, password string) (bool, uint, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "ValidateCredentials")
	defer span.End()
	user, err := service.GetUserInfo(ctx, username)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return false, 0, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return false, 0, err
	}
	customlogs.OtelLogger.Info(fmt.Sprintf("%s logged in!!", username))
	span.SetStatus(codes.Ok, "logged in")
	span.SetAttributes(attribute.Int64("userID", int64(user.ID)))
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

func (service *LoginServiceImpl) RegisterUser(ctx context.Context, username, password string) (models.Users, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "RegisterUser")
	defer span.End()
	_, err := service.GetUserInfo(ctx, username)
	if err == nil {
		customlogs.OtelLogger.Error(fmt.Sprintf("username: %s already exists", username))
		span.SetStatus(codes.Error, "username already exists")
		return models.Users{}, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Users{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return models.Users{}, err
	}

	secretKey, err := generateRandomString(16)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return models.Users{}, err
	}

	user := models.Users{Username: username, Password: string(hashedPassword), SecretKey: secretKey, Role: "user"}
	if err := service.DB.WithContext(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			customlogs.OtelLogger.Error(fmt.Sprintf("username: %s already exists", username))
			span.SetStatus(codes.Error, err.Error())
			return models.Users{}, errors.New("user with same username exists")
		}
		span.SetStatus(codes.Error, err.Error())
		return models.Users{}, err
	}
	customlogs.OtelLogger.Info(fmt.Sprintf("User: %s registered successfully!", username))
	span.SetStatus(codes.Ok, "user registered")
	span.SetAttributes(attribute.Int64("userID", int64(user.ID)))
	return user, nil
}
