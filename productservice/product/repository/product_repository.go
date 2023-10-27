package repository

import (
	"context"
	"errors"

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

func (r *productRepository) Store(ctx context.Context, product domain.Product) error {
	if err := r.db.Where("name", product.Name).First(&product).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("product name duplicated")
	}
	err := r.db.Debug().Create(&product).Error
	return err
}

func (r *productRepository) GetById(id uint32) (domain.Product, error) {
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
