package config

import (
	"albums/internal/database"
	"albums/internal/handlers"
	"albums/internal/middleware"
	"albums/internal/service"
	"albums/routes"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func initDB(envVars map[string]string) (*gorm.DB, error) {
	var db *gorm.DB

	if envVars["DB_TYPE"] == "local" {
		sqlite, err := database.SqliteConnect()
		if err != nil {
			return nil, err
		}
		err = database.SqliteInit(sqlite)
		if err != nil {
			return nil, err
		}
		db = sqlite
	} else {
		mysql, err := database.MysqlConnect(envVars["DB_HOST"], envVars["DB_USER"], envVars["DB_PASSWORD"], envVars["DB_NAME"])
		if err != nil {
			return nil, err
		}
		db = mysql
	}
	return db, nil
}
func RouterSetup() *gin.Engine {
	envVars, err := loadEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	db, err := initDB(envVars)
	if err != nil {
		log.Fatalln(err.Error())
	}
	router := gin.New()
	router.Use(middleware.LoggerMiddleware(), gin.Recovery())
	albumService := &service.AlbumServiceImpl{DB: db}
	albumHandler := &handlers.AlbumHandler{AlbumService: albumService}
	routes.RegisterAlbumRoutes(router, albumHandler)

	loginService := &service.LoginServiceImpl{DB: db}
	loginHandler := &handlers.LoginHandler{LoginService: loginService}
	routes.RegisterLoginRoutes(router, loginHandler, loginService)
	return router
}
