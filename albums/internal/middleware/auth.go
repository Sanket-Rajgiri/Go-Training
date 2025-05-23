package middleware

import (
	"albums/internal/models"
	"albums/internal/service"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var BasicAuthAccounts = gin.Accounts{
	"user": "admin",
}

var RoleMapping = map[string][]string{
	"user": []string{
		"POST /login",
		"POST /login/Register",
		"GET /albums/",
		"GET /albums/:id",
		"POST /albums/",
		"PATCH /albums/"},
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		bearerToken := strings.Split(authHeader, " ")[1]
		if _, ok := service.TokenStore[bearerToken]; !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		}
		c.Next()
	}
}
func BasicAuthMiddleware() gin.HandlerFunc {
	return gin.BasicAuth(BasicAuthAccounts)
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Empty token"})
			return
		}
		claims := &service.TokenClaim{}
		_, err := jwt.ParseWithClaims(bearerToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("Role", claims.Role)
		c.Next()
	}
}

func RoleValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("Role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
			return
		}
		role, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid role type"})
			return
		}
		allowedRoutes, ok := RoleMapping[role]
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
			return
		}
		route := c.Request.Method + " " + c.FullPath()
		for _, expectedRoute := range allowedRoutes {
			if route == expectedRoute {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "route not allowed to access", "route": route})
	}
}

func DBAuthMiddleware(service service.LoginService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials models.LoginPayload
		if err := c.ShouldBindBodyWithJSON(&credentials); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json", "msg": err.Error()})
			return
		}
		authenticated, userID, err := service.ValidateCredentials(c.Request.Context(), credentials.Username, credentials.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !authenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
			return
		}
		c.Set("UserID", strconv.FormatUint(uint64(userID), 10))
		c.Next()
	}
}
