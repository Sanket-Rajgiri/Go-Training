package database

import (
	"albums/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var albums = []models.Album{
		{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
	if DB.Migrator().HasTable(&models.Album{}) {
		if err := DB.Migrator().DropTable(&models.Album{}); err != nil {
			log.Printf("Error Dropping Table: %s", err)
			return
		}
	}
	if err := DB.Migrator().CreateTable(&models.Album{}); err != nil {
		log.Printf("Error Creating Table: %s", err)
		return
	}
	if err := DB.Create(&albums).Error; err != nil {
		log.Printf("Error Initialising DB: %s", err)
		return
	}
	log.Println("Initialised Successfully.")
}

func Connect() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error Connecting DB : %s", err)
		return
	}
	DB = db
	log.Println("Connected to Database")
}
