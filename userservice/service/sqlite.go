package service

import (
	"github.com/funthere/starset/userservice/domain"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("db1.db"), &gorm.Config{})

	// db.Debug().Migrator().DropTable(domain.User{})

	db.Debug().AutoMigrate(domain.User{})

	return db, err
}
