package handlers

import (
	"albums/internal/models"
	"albums/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	AlbumService service.AlbumsService
}

// GetAlbums godoc
//
//	@Summary		Get all albums
//	@Description	Returns all albums with ID, title, artist, and price
//	@Tags			albums
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string][]models.AlbumSwagger
//	@Failure		500	{object}	map[string]string
//	@Router			/albums [get]
//	@Security		BearerAuth
func (handler *AlbumHandler) GetAlbums(c *gin.Context) {
	albums, err := handler.AlbumService.GetAlbums(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"albums": albums})
}

// GetAlbumByID godoc
//
//	@Summary		Get album by ID
//	@Description	Returns a single album by its ID
//	@Tags			albums
//	@Produce		json
//	@Param			id	path		int	true	"Album ID"
//	@Success		200	{object}	models.AlbumSwagger
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/albums/{id} [get]
//
//	@Security		BearerAuth
func (handler *AlbumHandler) GetAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	album, err := handler.AlbumService.GetAlbumByID(c.Request.Context(), uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, album)

}

// AddAlbums godoc
//
//	@Summary		Add a new album
//	@Description	Creates a new album and returns its ID
//	@Tags			albums
//	@Accept			json
//	@Produce		json
//	@Param			album	body		models.AlbumSwagger	true	"Album to add"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/albums [post]
//
//	@Security		BearerAuth
func (handler *AlbumHandler) AddAlbums(c *gin.Context) {
	var newAlbum models.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	id, err := handler.AlbumService.AddAlbums(c.Request.Context(), newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Album Added", "ID": id})
}

// UpdatePrice godoc
//
//	@Summary		Update album price
//	@Description	Updates the price of an existing album
//	@Tags			albums
//	@Accept			json
//	@Produce		json
//	@Param			album	body		models.AlbumSwagger	true	"Album with updated price"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/albums [put]
//
//	@Security		BearerAuth
func (handler *AlbumHandler) UpdatePrice(c *gin.Context) {
	var updateAlbum models.Album
	if err := c.BindJSON(&updateAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if _, err := handler.AlbumService.UpdatePrice(c.Request.Context(), updateAlbum); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Album Price Updated Successfully"})
}

// DeleteAlbum godoc
//
//	@Summary		Delete album
//	@Description	Deletes an album by ID
//	@Tags			albums
//	@Produce		json
//	@Param			id	path		int	true	"Album ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/albums/{id} [delete]
//
//	@Security		BearerAuth
func (handler *AlbumHandler) DeleteAlbum(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	deletedAlbumID, err := handler.AlbumService.DeleteAlbum(c.Request.Context(), uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "album deleted successfully", "ID": deletedAlbumID})
}
