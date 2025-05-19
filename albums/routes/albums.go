package routes

import (
	"albums/internal/handlers"
	"albums/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAlbumRoutes(router *gin.Engine, handler *handlers.AlbumHandler) {
	albums := router.Group("/albums", middleware.JwtAuthMiddleware(), middleware.RoleValidationMiddleware())
	albums.GET("/", handler.GetAlbums)
	albums.GET("/:id", handler.GetAlbumByID)
	albums.POST("/", handler.AddAlbums)
	albums.PATCH("/", handler.UpdatePrice)
	albums.DELETE("/:id", handler.DeleteAlbum)
}
