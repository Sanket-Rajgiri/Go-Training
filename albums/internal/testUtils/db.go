package testutils

import (
	"albums/internal/models"
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitTestDB(ctx context.Context) (*gorm.DB, func(), error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to in-memory DB: %w", err)
	}

	// Auto-migrate your models
	if err := db.AutoMigrate(&models.Album{}, &models.Users{}); err != nil {
		return nil, nil, fmt.Errorf("migration failed: %w", err)
	}

	// Cleanup function to close DB (not strictly needed for SQLite, but good for MySQL/Postgres)
	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return db, cleanup, nil
}
