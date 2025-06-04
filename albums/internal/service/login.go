package service

import (
	"albums/internal/config/env"
	"albums/internal/customlogs"
	"albums/internal/models"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrUserAlreadyExists  = errors.New("username already exists")
	ErrInvalidTokenClaims = errors.New("invalid token claim")
)

type TokenClaim struct {
	Role string `json:"Role"`
	jwt.RegisteredClaims
}

type LoginServiceImpl struct {
	DB *gorm.DB
}

func (service *LoginServiceImpl) Login(ctx context.Context, userID, role string) (string, string, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "Login")
	defer span.End()
	var access_token, refresh_token string
	user, err := service.getUserbyID(ctx, userID)
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error fetching user details: %v", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return access_token, refresh_token, err
	}

	access_token, err = service.jwtTokenGenerator(ctx, userID, role, []byte(env.JWT_SECRET), env.JWT_EXPIRY_MINUTES)
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error generating access token: %v", err.Error()))

		span.SetStatus(codes.Error, err.Error())
		return access_token, refresh_token, err
	}

	refresh_token, err = service.jwtTokenGenerator(ctx, userID, role, []byte(user.SecretKey), env.REFRESH_TOKEN_EXPIRY_MINUTES)
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error generating refresh token: %v", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return access_token, refresh_token, err

	}

	user.Revoked = false
	if err := service.DB.WithContext(ctx).Save(user).Error; err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error updating user: %v", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return access_token, refresh_token, err
	}

	return access_token, refresh_token, nil
}

func (service *LoginServiceImpl) Logout(ctx context.Context, username string) error {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "Logout")
	defer span.End()
	user, err := service.getUserInfo(ctx, username)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
		return err
	}
	user.Revoked = true
	if err := service.DB.WithContext(ctx).Save(user).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
		return err
	}
	customlogs.OtelLogger.Ctx(ctx).Info(fmt.Sprintf("User : %s logged out", username))
	span.SetStatus(codes.Ok, "logged out")
	span.SetAttributes(attribute.Int64("userID", int64(user.ID)))
	return nil
}

func (service *LoginServiceImpl) RefreshToken(ctx context.Context, username, token string) (string, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "RefreshToken")
	defer span.End()
	var newAccessToken string
	user, err := service.getUserInfo(ctx, username)
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error fetching user details: %v", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return newAccessToken, err
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	claims, err := service.jwtTokenValidator(ctx, token, user.SecretKey)
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error validating token: %v", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return newAccessToken, err
	}

	if claims.Role != user.Role || claims.RegisteredClaims.Subject != userID {
		customlogs.OtelLogger.Ctx(ctx).Error(ErrInvalidTokenClaims.Error())
		span.SetStatus(codes.Error, ErrInvalidTokenClaims.Error())
		return newAccessToken, ErrInvalidTokenClaims
	}

	newAccessToken, err = service.jwtTokenGenerator(ctx, userID, user.Role, []byte(env.JWT_SECRET), env.JWT_EXPIRY_MINUTES)
	if err != nil {
		customlogs.OtelLogger.Error(fmt.Sprintf("error generating access token: %s", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return newAccessToken, err
	}
	customlogs.OtelLogger.Ctx(ctx).Info(fmt.Sprintf("user: %s token refreshed", username))
	span.SetStatus(codes.Ok, "token generated")
	span.SetAttributes(attribute.String("userID", userID))
	return newAccessToken, nil
}

func (service *LoginServiceImpl) RegisterUser(ctx context.Context, username, password string) (models.Users, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "RegisterUser")
	defer span.End()
	_, err := service.getUserInfo(ctx, username)
	if err == nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("username: %s already exists", username))
		span.SetStatus(codes.Error, ErrUserAlreadyExists.Error())
		return models.Users{}, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Users{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
		return models.Users{}, err
	}

	secretKey, err := service.generateRandomString(16)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
		return models.Users{}, err
	}

	user := models.Users{Username: username, Password: string(hashedPassword), SecretKey: secretKey, Role: "user"}
	if err := service.DB.WithContext(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("username: %s already exists", username))
			span.SetStatus(codes.Error, ErrUserAlreadyExists.Error())
			return models.Users{}, ErrUserAlreadyExists
		}
		span.SetStatus(codes.Error, err.Error())
		return models.Users{}, err
	}
	customlogs.OtelLogger.Ctx(ctx).Info(fmt.Sprintf("User: %s registered successfully!", username))
	span.SetStatus(codes.Ok, "user registered")
	span.SetAttributes(attribute.Int64("userID", int64(user.ID)))
	return user, nil
}

func (service *LoginServiceImpl) ValidateCredentials(ctx context.Context, username, password string) (uint, string, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "ValidateCredentials")
	defer span.End()
	user, err := service.getUserInfo(ctx, username)
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error fetching user details: %v", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return 0, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error validating password: %v", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return 0, "", err
	}
	customlogs.OtelLogger.Ctx(ctx).Info("credentials valid")
	span.SetStatus(codes.Ok, "credentials valid")
	span.SetAttributes(attribute.Int64("userID", int64(user.ID)))
	return user.ID, user.Role, nil
}

func (service *LoginServiceImpl) generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	var randomString string
	_, err := rand.Read(bytes)
	if err != nil {
		return randomString, err
	}

	for i, b := range bytes {
		bytes[i] = charset[int(b)%len(charset)]
	}
	randomString = string(bytes)
	return randomString, nil
}

func (service *LoginServiceImpl) jwtTokenValidator(ctx context.Context, bearerToken, secretKey string) (*TokenClaim, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "jwtTokenValidator")
	defer span.End()
	claims := &TokenClaim{}
	_, err := jwt.ParseWithClaims(bearerToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	span.SetStatus(codes.Ok, "token validated")
	customlogs.OtelLogger.Ctx(ctx).Info("token validated")
	return claims, nil
}

func (service *LoginServiceImpl) jwtTokenGenerator(ctx context.Context, userID, role string, jwtSecret []byte, jwtExpiryMinutes int) (string, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "jwtTokenGenerator")
	defer span.End()
	var signedToken string
	expiryTime := time.Now().Add(time.Minute * time.Duration(jwtExpiryMinutes))
	claims := TokenClaim{
		Role: role,
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
		customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
		return signedToken, err
	}
	span.SetStatus(codes.Ok, "token generated")
	return signedToken, nil
}

func (service *LoginServiceImpl) getUserbyID(ctx context.Context, userID string) (models.Users, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "getUserbyID")
	defer span.End()
	var user models.Users
	if err := service.DB.WithContext(ctx).Find(&user, userID).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return models.Users{}, err
	}
	span.SetStatus(codes.Ok, "user fetched")
	span.SetAttributes(attribute.Int64("userID", int64(user.ID)))
	return user, nil
}

func (service *LoginServiceImpl) getUserInfo(ctx context.Context, username string) (models.Users, error) {
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
