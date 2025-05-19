package routes

import (
	"albums/handlers"
	"albums/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAlbumRoutes(router *gin.Engine) {
	albums := router.Group("/albums", middleware.JwtAuthMiddleware(), middleware.RoleValidationMiddleware())
	albums.GET("/", handlers.GetAlbums)
	albums.GET("/:id", handlers.GetAlbumByID)
	albums.POST("/", handlers.AddAlbums)
	albums.PATCH("/", handlers.UpdatePrice)
	albums.DELETE("/:id", handlers.DeleteAlbum)
}
