package service

import (
	"github.com/funthere/starset/orderservice/domain"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("order.db"), &gorm.Config{})

	// db.Debug().Migrator().DropTable(domain.Order{})

	db.Debug().AutoMigrate(domain.Order{})
	db.Debug().AutoMigrate(domain.OrderProduct{})

	return db, err
}
