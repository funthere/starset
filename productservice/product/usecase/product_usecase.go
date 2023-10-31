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
	if err := u.productRepo.Store(ctx, product); err != nil {
		return err
	}
	ids := []uint32{product.OwnerID}
	mapUsers, err := u.userService.GetUserByIds(ctx, ids)
	if err != nil {
		return nil
	}

	if len(mapUsers) > 0 {
		if user, ok := mapUsers[product.OwnerID]; ok {
			product.Owner.ID = user.ID
			product.Owner.Name = user.Name
			product.Owner.Email = user.Email
			product.Owner.Address = user.Address
		}
	}

	return nil
}

func (u *productUsecase) GetById(ctx context.Context, id uint32) (domain.Product, error) {
	product, err := u.productRepo.GetById(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}
	ids := []uint32{product.OwnerID}

	mapUsers, err := u.userService.GetUserByIds(ctx, ids)
	if err != nil {
		return product, nil
	}

	if len(mapUsers) > 0 {
		if user, ok := mapUsers[product.OwnerID]; ok {
			product.Owner.ID = user.ID
			product.Owner.Name = user.Name
			product.Owner.Email = user.Email
			product.Owner.Address = user.Address
		}
	}

	return product, nil
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
		return producs, nil
	}

	// fill the product.owner attributes
	if len(mapUsers) > 0 {
		for i := range producs {
			if user, ok := mapUsers[producs[i].OwnerID]; ok {
				producs[i].Owner.Name = user.Name
				producs[i].Owner.Email = user.Email
				producs[i].Owner.Address = user.Address
			}
		}
	}

	return producs, nil
}
