package middleware

import (
	"albums/internal/config/env"
	"albums/internal/customlogs"
	"albums/internal/handlers"
	"albums/internal/models"
	"albums/internal/service"

	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var (
	BasicAuthAccounts = gin.Accounts{
		"user": "admin",
	}

	RoleMapping = map[string][]string{
		"user": []string{
			"GET /albums/",
			"GET /albums/:id",
			"POST /albums/",
			"PATCH /albums/",
		},
		"reader": []string{
			"GET /albums/",
			"GET /albums/:id",
		},
		"admin": []string{
			"GET /albums/",
			"GET /albums/:id",
			"POST /albums/",
			"PATCH /albums/",
			"DELETE /albums/:id",
		},
	}
)

func BasicAuthMiddleware() gin.HandlerFunc {
	return gin.BasicAuth(BasicAuthAccounts)
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("Middleware-Tracer")
		ctx, span := tracer.Start(c.Request.Context(), "JwtAuthMiddleware")
		c.Request = c.Request.WithContext(ctx)

		defer span.End()
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			customlogs.OtelLogger.Ctx(ctx).Error("Missing or invalid Authorization header")
			span.SetStatus(codes.Error, "Missing or invalid Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if bearerToken == "" {
			customlogs.OtelLogger.Ctx(ctx).Error("Empty token")
			span.SetStatus(codes.Error, "Empty Token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Empty token"})
			return
		}
		claims := &service.TokenClaim{}
		_, err := jwt.ParseWithClaims(bearerToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(env.JWT_SECRET), nil
		})
		if err != nil {
			customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("error in validating token: %v", err))
			span.SetStatus(codes.Error, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("Role", claims.Role)
		span.SetAttributes(attribute.String("role", claims.Role))
		span.SetStatus(codes.Ok, "token validated")
		c.Next()
	}
}

func RoleValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("Middleware-Tracer")
		ctx, span := tracer.Start(c.Request.Context(), "RoleValidationMiddleware")
		c.Request = c.Request.WithContext(ctx)
		defer span.End()
		roleVal, exists := c.Get("Role")
		if !exists {
			span.SetStatus(codes.Error, "role not found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
			return
		}
		role, ok := roleVal.(string)
		if !ok {
			span.SetStatus(codes.Error, "invalid role type")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid role type"})
			return
		}
		span.SetAttributes(attribute.String("role", role))
		allowedRoutes, ok := RoleMapping[role]
		if !ok {
			span.SetStatus(codes.Error, "Access Denied")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
			return
		}
		route := c.Request.Method + " " + c.FullPath()
		for _, expectedRoute := range allowedRoutes {
			if route == expectedRoute {
				span.SetStatus(codes.Ok, "access granted")
				c.Next()
				return
			}
		}
		span.SetStatus(codes.Error, "route not allowed to access")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "route not allowed to access", "route": route})
	}
}

func DBAuthMiddleware(service handlers.LoginService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("Middleware-Tracer")
		ctx, span := tracer.Start(c.Request.Context(), "DBAuthMiddleware")
		c.Request = c.Request.WithContext(ctx)
		defer span.End()
		var credentials models.LoginPayload
		if err := c.ShouldBindBodyWithJSON(&credentials); err != nil {
			customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
			span.SetStatus(codes.Error, err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json", "msg": err.Error()})
			return
		}
		userID, role, err := service.ValidateCredentials(c.Request.Context(), credentials.Username, credentials.Password)
		if err != nil {
			customlogs.OtelLogger.Ctx(ctx).Error(err.Error())
			span.SetStatus(codes.Error, err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Set("UserID", strconv.FormatUint(uint64(userID), 10))
		c.Set("Role", role)
		span.SetStatus(codes.Ok, "authenticated user")
		c.Next()
	}
}
