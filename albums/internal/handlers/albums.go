package handlers

import (
	"albums/internal/models"
	"albums/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("Albums-Tracer")
	ctx, span := tracer.Start(ctx, "GetAlbumsHandler")
	defer span.End()
	albums, err := handler.AlbumService.GetAlbums(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	span.SetStatus(codes.Ok, "albums fetched")
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("Login-Tracer")
	_, span := tracer.Start(ctx, "GetAlbumByIDHandler")
	defer span.End()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	album, err := handler.AlbumService.GetAlbumByID(ctx, uint(id))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	span.SetAttributes(attribute.Int64("albumId", int64(album.ID)))
	span.SetStatus(codes.Ok, "album fetched")
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("Albums-Tracer")
	_, span := tracer.Start(ctx, "AddAlbumsHandler")
	defer span.End()
	var newAlbum models.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	id, err := handler.AlbumService.AddAlbums(ctx, newAlbum)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	span.SetAttributes(attribute.Int64("albumId", int64(newAlbum.ID)))
	span.SetStatus(codes.Ok, "album added")
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("Albums-Tracer")
	_, span := tracer.Start(ctx, "UpdatePriceHandler")
	defer span.End()
	var updateAlbum models.Album
	if err := c.BindJSON(&updateAlbum); err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if _, err := handler.AlbumService.UpdatePrice(ctx, updateAlbum); err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	span.SetAttributes(attribute.Int64("albumId", int64(updateAlbum.ID)))
	span.SetStatus(codes.Ok, "album updated")
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("Albums-Tracer")
	_, span := tracer.Start(ctx, "DeleteAlbumHandler")
	defer span.End()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	deletedAlbumID, err := handler.AlbumService.DeleteAlbum(ctx, uint(id))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	span.SetAttributes(attribute.Int64("albumId", int64(deletedAlbumID)))
	span.SetStatus(codes.Ok, "album deleted")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "album deleted successfully", "ID": deletedAlbumID})
}
