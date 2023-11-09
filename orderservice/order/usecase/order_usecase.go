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
	var productIDs []int64
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
	orders, err := u.orderRepo.Fetch(ctx, filter)
	if err != nil {
		return []domain.Order{}, err
	}

	// Populate product data
	var productIDs []int64
	for i := range orders {
		for j := range orders[i].Items {
			productIDs = append(productIDs, orders[i].Items[j].ProductID)
		}
	}

	mapProducts, err := u.productService.GetProductByIds(ctx, productIDs)
	if err != nil {
		return []domain.Order{}, errors.New(fmt.Sprintf("order.usecase.fetch: %v", err.Error()))
	}

	for i := range orders {
		orders[i].Buyer.ID = orders[i].BuyerID
		orders[i].Seller.ID = orders[i].SellerID
		for j := range orders[i].Items {
			if product, ok := mapProducts[orders[i].Items[j].ProductID]; ok {
				orders[i].Items[j].Product.ID = product.ID
				orders[i].Items[j].Product.Name = product.Name
				orders[i].Items[j].Product.Description = product.Description
				orders[i].Items[j].Product.Price = product.Price
			}
		}
	}

	// Populate data of buyer and seller
	var userIDs []int64
	for i := range orders {
		userIDs = append(userIDs, orders[i].BuyerID)
		userIDs = append(userIDs, orders[i].SellerID)
	}

	mapUser, err := u.userService.GetUserByIds(ctx, userIDs)
	if err != nil {
		return orders, nil
	}
	for i := range orders {
		if user, ok := mapUser[orders[i].BuyerID]; ok {
			orders[i].Buyer = user
		}
		if user, ok := mapUser[orders[i].SellerID]; ok {
			orders[i].Seller = user
		}
	}

	return orders, nil
}

func (u orderUsecase) PatchStatus(ctx context.Context, id int64, status string) error {
	return u.orderRepo.PatchStatus(ctx, id, status)
}
