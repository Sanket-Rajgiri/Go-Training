package config

import (
	"albums/internal/config/env"
	"albums/internal/database"

	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	var db *gorm.DB

	if env.DB_TYPE == "local" {
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
		mysql, err := database.MysqlConnect(env.DB_HOST, env.DB_USER, env.DB_PASSWORD, env.DB_NAME)
		if err != nil {
			return nil, err
		}
		db = mysql
	}
	return db, nil
}
