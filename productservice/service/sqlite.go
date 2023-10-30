package service

import (
	"github.com/funthere/starset/productservice/domain"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("product.db"), &gorm.Config{})

	// db.Debug().Migrator().DropTable(domain.Product{})

	db.Debug().AutoMigrate(domain.Product{})

	return db, err
}
