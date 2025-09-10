package database

import (
	"fmt"
	"log"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlConnect(dbHost, dbUser, dbPassword, dbName string) (*gorm.DB, error) {
	mysql_dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(mysql_dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting db : %s", err)

	}
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("error initialising otelgorm: %s", err)
	}
	log.Println("Connected to Database")
	return db, nil
}
