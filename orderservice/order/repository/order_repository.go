package repository

import (
	"context"

	"github.com/funthere/starset/orderservice/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r orderRepository) Store(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Debug().Create(&order).Error
}

func (r orderRepository) Fetch(ctx context.Context, filter domain.Filter) ([]domain.Order, error) {
	orders := []domain.Order{}
	qBuilder := r.db.WithContext(ctx).Debug()
	if filter.BuyerID > 0 {
		qBuilder = qBuilder.Where("buyer_id", filter.BuyerID)
	}
	if filter.OrderID > 0 {
		qBuilder = qBuilder.Where("buyer_id", filter.OrderID)
	}
	if filter.SellerID > 0 {
		qBuilder = qBuilder.Where("seller_id", filter.SellerID)
	}

	if err := qBuilder.Find(&orders).Error; err != nil {
		return []domain.Order{}, err
	}

	return orders, nil
}

func (r orderRepository) PatchStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Debug().Model(&domain.Order{}).Where("id", id).Update("status", status).Error
}
