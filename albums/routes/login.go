package routes

import (
	"albums/internal/handlers"
	"albums/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(router *gin.Engine, handler *handlers.LoginHandler, service handlers.LoginService) {
	router.POST("/logout", middleware.JwtAuthMiddleware(), handler.Logout)
	login := router.Group("/login")
	// login.GET("/", handlers.LoginHandler, middleware.BasicAuthMiddleware())
	login.POST("/refresh", handler.Refresh)
	login.POST("/", middleware.DBAuthMiddleware(service), handler.Login)
	login.POST("/register", handler.Register)
}
