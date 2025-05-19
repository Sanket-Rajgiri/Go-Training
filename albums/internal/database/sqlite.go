package database

import (
	"albums/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SqliteInit(DB *gorm.DB) error {
	var albums = []models.Album{
		{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
	if DB.Migrator().HasTable(&models.Album{}) {
		if err := DB.Migrator().DropTable(&models.Album{}); err != nil {
			return fmt.Errorf("error dropping table: %s", err)
		}
	}
	if err := DB.Migrator().CreateTable(&models.Album{}); err != nil {
		return fmt.Errorf("Error Creating Table: %s", err)
	}
	if err := DB.Create(&albums).Error; err != nil {
		return fmt.Errorf("Error Initialising DB: %s", err)

	}
	log.Println("Initialised Successfully.")
	return nil
}

func SqliteConnect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Error Connecting DB : %s", err)
	}
	log.Println("Connected to Database")
	return db, nil
}
