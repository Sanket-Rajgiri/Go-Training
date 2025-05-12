package handlers

import (
	"albums/database"
	"albums/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// type album struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

func GetAlbums(c *gin.Context) {
	var albums []models.Album
	if err := database.DB.Find(&albums).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch albums"})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"albums": albums})
}
func GetAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var album models.Album
	if err := database.DB.First(&album, id); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"album": album})

}

func AddAlbums(c *gin.Context) {
	var newAlbum models.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if err := database.DB.Create(&newAlbum).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not add album",
			"message": err})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Album Added", "ID": newAlbum.ID})
}

func UpdatePrice(c *gin.Context) {
	var existingAlbum, updateAlbum models.Album
	if err := c.BindJSON(&updateAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if err := database.DB.First(&existingAlbum, updateAlbum.ID, updateAlbum.Artist, updateAlbum.Title).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}
	existingAlbum.Price = updateAlbum.Price
	if err := database.DB.Save(&existingAlbum).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not update album",
			"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Album Price Updated Successfully"})
}

func DeleteAlbum(c *gin.Context) {
	var album models.Album
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := database.DB.First(&album, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}
	if err := database.DB.Delete(&album).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "album deleted successfully"})
}
