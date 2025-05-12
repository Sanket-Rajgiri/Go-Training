package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func GetAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
func GetAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, album := range albums {
		if album.ID == id {
			c.IndentedJSON(http.StatusFound, album)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func AddAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Album Added", "ID": newAlbum.ID})
}

func UpdatePrice(c *gin.Context) {
	var updateAlbum album
	if err := c.BindJSON(&updateAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	for i, album := range albums {
		if album.ID == updateAlbum.ID {
			if album.Artist == updateAlbum.Artist && album.Title == updateAlbum.Title {
				albums[i] = updateAlbum
				c.IndentedJSON(http.StatusOK, gin.H{"message": "Price Updated for Album", "ID": updateAlbum.ID})
				return
			} else {
				c.IndentedJSON(http.StatusPreconditionFailed, gin.H{"error": "Album's details does not match"})
				return
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	var sliceIndex = -1

	for index, album := range albums {
		if album.ID == id {
			sliceIndex = index
			break
		}
	}
	if sliceIndex < 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No album exists with ID %s", id)})
		return
	}
	albums = append(albums[:sliceIndex], albums[sliceIndex+1:]...)
	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Album deleted with ID %s", id)})
}
