package main

import (
	"albums/internal/config"
	"log"
)

//	@title						Go Gin Rest API
//	@version					1.0
//	@description				A rest API in Go using Gin framework
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your bearer token in the format **Bearer &lt;token&gt;**

func main() {

	router := config.RouterSetup()
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error Starting Server : %v", err)
	}
}
