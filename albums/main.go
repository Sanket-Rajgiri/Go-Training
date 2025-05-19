package main

import (
	"albums/internal/config"
	"log"
)

func main() {

	router := config.RouterSetup()
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error Starting Server : %v", err)
	}
}
