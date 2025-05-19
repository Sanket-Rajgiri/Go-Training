package handlers

import (
	"albums/internal/models"
	"albums/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//	type album struct {
//		ID     string  `json:"id"`
//		Title  string  `json:"title"`
//		Artist string  `json:"artist"`
//		Price  float64 `json:"price"`
//	}
type AlbumHandler struct {
	AlbumService service.AlbumsService
}

func (handler *AlbumHandler) GetAlbums(c *gin.Context) {
	albums, err := handler.AlbumService.GetAlbums()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"albums": albums})
}
func (handler *AlbumHandler) GetAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	album, err := handler.AlbumService.GetAlbumByID(uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"album": album})

}

func (handler *AlbumHandler) AddAlbums(c *gin.Context) {
	var newAlbum models.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	id, err := handler.AlbumService.AddAlbums(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Album Added", "ID": id})
}

func (handler *AlbumHandler) UpdatePrice(c *gin.Context) {
	var updateAlbum models.Album
	if err := c.BindJSON(&updateAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if _, err := handler.AlbumService.UpdatePrice(updateAlbum); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Album Price Updated Successfully"})
}

func (handler *AlbumHandler) DeleteAlbum(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	deletedAlbumID, err := handler.AlbumService.DeleteAlbum(uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "album deleted successfully", "ID": deletedAlbumID})
}
