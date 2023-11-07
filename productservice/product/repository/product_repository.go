package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/funthere/starset/productservice/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Store(ctx context.Context, product *domain.Product) error {
	if err := r.db.Where("name", product.Name).First(product).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("product name duplicated")
	}
	err := r.db.WithContext(ctx).Debug().Create(&product).Error
	return err
}

func (r *productRepository) GetById(ctx context.Context, id uint32) (domain.Product, error) {
	product := domain.Product{}
	if err := r.db.First(&product, id).Error; err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (r *productRepository) FetchByIds(ctx context.Context, ids []uint32) ([]domain.Product, error) {
	products := []domain.Product{}
	if err := r.db.Find(&products, ids).Error; err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}

func (r *productRepository) Fetch(ctx context.Context, filter domain.Filter) ([]domain.Product, error) {
	products := []domain.Product{}
	qBuilder := r.db.WithContext(ctx).Debug()
	if filter.OrderID > 0 {
		qBuilder = qBuilder.Where("id", filter.OrderID)
	}

	if filter.OwnerID > 0 {
		qBuilder = qBuilder.Where("user_id", filter.OwnerID)
	}

	if filter.Search != "" {
		qBuilder = qBuilder.Where("name LIKE ?", fmt.Sprintf("%%%v%%", filter.Search))
	}

	if len(filter.IDs) > 0 {
		qBuilder = qBuilder.Where("id", filter.IDs)
	}

	if err := qBuilder.Find(&products).Error; err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}
