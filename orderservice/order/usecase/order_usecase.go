package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/funthere/starset/orderservice/domain"
	"github.com/funthere/starset/orderservice/service/product"
	"github.com/funthere/starset/orderservice/service/user"
)

type orderUsecase struct {
	orderRepo      domain.OrderRepository
	productService product.ProductService
	userService    user.UserService
}

func NewOrderUsecase(
	repo domain.OrderRepository,
	productSvc product.ProductService,
	userSvc user.UserService,
) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:      repo,
		productService: productSvc,
		userService:    userSvc,
	}
}

func (u orderUsecase) Store(ctx context.Context, order *domain.Order) error {
	productIDs := []int64{}
	for i := range order.Items {
		productIDs = append(productIDs, order.Items[i].ProductID)
	}

	mapProducts, err := u.productService.GetProductByIds(ctx, productIDs)
	if err != nil {
		return errors.New(fmt.Sprintf("order.usecase.store: %v", err.Error()))
	}

	for i := range order.Items {
		order.Items[i].Product.ID = order.Items[i].ProductID
		if product, ok := mapProducts[order.Items[i].ProductID]; ok {
			order.Items[i].Product.Name = product.Name
			order.Items[i].Product.Description = product.Description
			order.Items[i].Product.Price = product.Price

			order.SellerID = product.Owner.ID
			order.Seller = product.Owner
			order.TotalPrice += order.Items[i].Product.Price * order.Items[i].Quantity
		} else {
			return errors.New(fmt.Sprintf("order.usecase.store: product with ID %v not found", order.Items[i].ProductID))
		}
	}

	return u.orderRepo.Store(ctx, order)
}

func (u orderUsecase) Fetch(ctx context.Context, filter domain.Filter) ([]domain.Order, error) {
	return u.orderRepo.Fetch(ctx, filter)
}

func (u orderUsecase) PatchStatus(ctx context.Context, id int64, status string) error {
	return u.PatchStatus(ctx, id, status)
}
