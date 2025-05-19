package routes

import (
	"albums/handlers"
	"albums/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(router *gin.Engine) {
	login := router.Group("/login")
	// login.GET("/", handlers.LoginHandler, middleware.BasicAuthMiddleware())
	login.POST("/", middleware.DBAuthMiddleware(), handlers.LoginHandler)
	login.POST("/register", handlers.RegisterHandler)
}
