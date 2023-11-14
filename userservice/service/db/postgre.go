package service

import (
	"fmt"
	"os"

	"github.com/funthere/starset/userservice/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgreDatabase() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.Debug().AutoMigrate(domain.User{})

	return db, err
}
