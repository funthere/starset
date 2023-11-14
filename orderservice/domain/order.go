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
	SourceAddress      string         `json:"source_address" gorm:"NULL;type:varchar(255);"`
	DestinationAddress string         `json:"destination_address" gorm:"NOT NULL;type:varchar(255);" validate:"required"`
	Notes              string         `json:"notes" gorm:"NULL;type:varchar(255);" `
	Items              []OrderProduct `json:"items" gorm:"foreignKey:OrderID" validate:"required,min=1"`
	TotalPrice         int64          `json:"total_price" gorm:"column:total_price;NOT NULL;"`
	Status             string         `json:"status" gorm:"NOT NULL;"`
	BuyerID            int64          `json:"-" gorm:"column:buyer_id;NULL;" validate:"-"`
	Buyer              User           `json:"buyer" gorm:"-"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type OrderProduct struct {
	ID        int64   `json:"-" gorm:"primaryKey"`
	Product   Product `json:"product" gorm:"-"`
	OrderID   int64   `json:"-" gorm:"column:order_id;NOT NULL;"`
	ProductID int64   `json:"product_id" gorm:"column:product_id;NOT NULL;" validate:"required"`
	Quantity  int64   `json:"quantity" gorm:"column:quantity;NOT NULL;" validate:"min=1"`
}

type User struct {
	ID      int64  `json:"-"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type Product struct {
	ID          int64  `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Owner       User   `json:"-"`
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
