package domain

import (
	"context"
	"time"
)

const (
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
)

type Order struct {
	ID                 int64          `json:"-" gorm:"primaryKey"`
	BuyerID            int64          `json:"-" gorm:"column:buyer_id;NULL;type:integer;" validate:"required"`
	Buyer              User           `json:"buyer" gorm:"-"`
	SellerID           int64          `json:"-" gorm:"column:seller_id;NULL;type:integer;" validate:"required"`
	Seller             User           `json:"seller" gorm:"-"`
	SourceAddress      string         `json:"source_address" gorm:"NULL;type:varchar(255);"`
	DestinationAddress string         `json:"destination_address" gorm:"NOT NULL;type:varchar(255);" validate:"required"`
	Notes              string         `json:"notes" gorm:"NULL;type:varchar(255);" `
	Items              []OrderProduct `json:"items" gorm:"foreignKey:ID" validate:"required,min=1"`
	TotalPrice         int64          `json:"total_price" gorm:"column:total_price;NOT NULL;type:integer;"`
	Status             string         `json:"status" gorm:"NOT NULL;"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type OrderProduct struct {
	ID        int64   `json:"-" gorm:"primaryKey"`
	Product   Product `json:"product" gorm:"-"`
	ProductID int64   `json:"product_id" gorm:"column:product_id;NOT NULL;type:integer;" validate:"required"`
	Quantity  int64   `json:"quantity" gorm:"column:quantity;NOT NULL;type:integer;" validate:"min=1"`
}

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Owner       User   `json:"owner"`
}

type Filter struct {
	BuyerID,
	SellerID int64
	OrderID int64
	Search  string
}

type OrderUsecase interface {
	Store(ctx context.Context, order *Order) error
	Fetch(ctx context.Context, filter Filter) ([]Order, error)
	PatchStatus(ctx context.Context, id int64, status string) error
}

type OrderRepository interface {
	Store(ctx context.Context, order *Order) error
	Fetch(ctx context.Context, filter Filter) ([]Order, error)
	PatchStatus(ctx context.Context, id int64, status string) error
}
