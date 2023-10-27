package usecase

import (
	"context"

	"github.com/funthere/starset/productservice/domain"
)

type productUsecase struct {
	productRepo domain.ProductUsecase
}

func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) Store(ctx context.Context, product domain.Product) error {
	return u.productRepo.Store(ctx, product)
}

func (u *productUsecase) GetById(id uint32) (domain.Product, error) {
	return u.productRepo.GetById(id)
}

func (u *productUsecase) FetchByIds(ctx context.Context, ids []uint32) ([]domain.Product, error) {
	return u.productRepo.FetchByIds(ctx, ids)
}
