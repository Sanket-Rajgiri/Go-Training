package routes

import (
	"albums/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAlbumRoutes(routerGroup *gin.RouterGroup) {
	albums := routerGroup.Group("/albums")
	albums.GET("/", handlers.GetAlbums)
	albums.GET("/:id", handlers.GetAlbumByID)
	albums.POST("/", handlers.AddAlbums)
	albums.PATCH("/", handlers.UpdatePrice)
	albums.DELETE("/:id", handlers.DeleteAlbum)
}
