package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint32    `json:"-" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"NOT NULL;type:varchar(255);" validate:"required"`
	Description string    `json:"description" gorm:"NOT NULL;type:varchar(255);" validate:"required"`
	Price       int64     `json:"price" gorm:"NOT NULL;" validate:"required"`
	OwnerID     uint32    `json:"-" gorm:"column:user_id;NULL;type:integer;" validate:"-"`
	Owner       User      `json:"owner" gorm:"-" validate:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New().ID()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	return nil
}

type Filter struct {
	OwnerID uint32
	OrderID uint32
	Search  string
	IDs     []uint32
}

type User struct {
	ID      uint32 `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}
type ProductUsecase interface {
	Store(ctx context.Context, product *Product) error
	GetById(ctx context.Context, id uint32) (Product, error)
	FetchByIds(ctx context.Context, ids []uint32) ([]Product, error)
	Fetch(ctx context.Context, filter Filter) ([]Product, error)
}

type ProductRepository interface {
	Store(ctx context.Context, product *Product) error
	GetById(ctx context.Context, id uint32) (Product, error)
	FetchByIds(ctx context.Context, ids []uint32) ([]Product, error)
	Fetch(ctx context.Context, filter Filter) ([]Product, error)
}
