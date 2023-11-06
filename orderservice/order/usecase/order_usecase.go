package usecase

import (
	"context"

	"github.com/funthere/starset/orderservice/domain"
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
}

func NewOrderUsecase(repo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo: repo,
	}
}

func (u orderUsecase) Store(ctx context.Context, order *domain.Order) error {
	return u.orderRepo.Store(ctx, order)
}

func (u orderUsecase) Fetch(ctx context.Context, filter domain.Filter) ([]domain.Order, error) {
	return u.orderRepo.Fetch(ctx, filter)
}

func (u orderUsecase) PatchStatus(ctx context.Context, id int64, status string) error {
	return u.PatchStatus(ctx, id, status)
}
