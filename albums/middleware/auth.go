package middleware

import (
	"albums/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var BasicAuthAccounts = gin.Accounts{
	"user": "admin",
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
		authHeader := c.GetHeader("Authorization")
		bearerToken := strings.Split(authHeader, " ")[1]
		err := service.JWTTokenValidator(bearerToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		c.Next()
	}
}
