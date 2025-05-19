package routes

import (
	"albums/internal/handlers"
	"albums/internal/middleware"
	"albums/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(router *gin.Engine, handler *handlers.LoginHandler, service service.LoginService) {
	login := router.Group("/login")
	// login.GET("/", handlers.LoginHandler, middleware.BasicAuthMiddleware())
	login.POST("/", middleware.DBAuthMiddleware(service), handler.Login)
	login.POST("/register", handler.Register)
}
