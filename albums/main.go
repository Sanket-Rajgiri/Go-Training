package main

import (
	"albums/database"
	"albums/routes"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var dbTypeMap = map[string]string{
	"dev":   "mysql",
	"local": "sqlite",
}

func loadEnv() (map[string]string, error) {
	envVariables := make(map[string]string)
	envName, exists := os.LookupEnv("ENV")
	if !exists {
		envName = "local"
		log.Println("setting envName = local")
	}
	envVariables["DB_TYPE"] = dbTypeMap[envName]
	if envName != "local" {
		requiredVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_NAME"}
		for _, key := range requiredVars {
			value := os.Getenv(key)
			if value == "" {
				return nil, errors.New("missing required environment variable: " + key)
			}
			envVariables[key] = value
		}
	}

	return envVariables, nil
}

func main() {

	router := gin.Default()
	rootPath := router.Group("/")
	routes.RegisterAlbumRoutes(rootPath)
	envVars, err := loadEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	if envVars["DB_TYPE"] == "local" {
		database.SqliteConnect()
		database.SqliteInit()
	} else {
		database.MysqlConnect(envVars["DB_HOST"], envVars["DB_USER"], envVars["DB_PASSWORD"], envVars["DB_NAME"])
	}
	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error Starting Server : %v", err)
	}
}
