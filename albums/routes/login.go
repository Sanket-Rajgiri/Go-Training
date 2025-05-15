package routes

import (
	"albums/handlers"
	"albums/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(router *gin.Engine) {
	login := router.Group("/login", middleware.BasicAuthMiddleware())
	login.GET("/", handlers.LoginHandler)
}
