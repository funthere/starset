package usecase

import (
	"context"

	"github.com/funthere/starset/productservice/domain"
	"github.com/funthere/starset/productservice/service/user"
)

type productUsecase struct {
	productRepo domain.ProductUsecase
	userService user.UserService
}

func NewProductUsecase(productRepo domain.ProductRepository, userService user.UserService) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
		userService: userService,
	}
}

func (u *productUsecase) Store(ctx context.Context, product *domain.Product) error {
	return u.productRepo.Store(ctx, product)
}

func (u *productUsecase) GetById(id uint32) (domain.Product, error) {
	return u.productRepo.GetById(id)
}

func (u *productUsecase) FetchByIds(ctx context.Context, ids []uint32) ([]domain.Product, error) {
	return u.productRepo.FetchByIds(ctx, ids)
}

func (u *productUsecase) Fetch(ctx context.Context, filter domain.Filter) ([]domain.Product, error) {
	producs, err := u.productRepo.Fetch(ctx, filter)
	if err != nil {
		return []domain.Product{}, err
	}

	ids := []uint32{}
	for i := range producs {
		producs[i].Owner.ID = producs[i].OwnerID
		ids = append(ids, producs[i].OwnerID)
	}

	mapUsers, err := u.userService.GetUserByIds(ctx, ids)
	if err != nil {
		return []domain.Product{}, err
	}

	// fill the product.owner attributes
	for i := range producs {
		if _, ok := mapUsers[producs[i].OwnerID]; ok {
			producs[i].Owner.Name = mapUsers[producs[i].OwnerID].Name
			producs[i].Owner.Email = mapUsers[producs[i].OwnerID].Email
			producs[i].Owner.Address = mapUsers[producs[i].OwnerID].Address
		}
	}

	return producs, nil
}
