package main

import (
	"albums/database"
	"albums/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	rootPath := router.Group("/")
	routes.RegisterAlbumRoutes(rootPath)

	database.Connect()
	database.Init()
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error Starting Server : %v", err)
	}
}
