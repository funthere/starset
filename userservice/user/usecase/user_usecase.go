package usecase

import (
	"context"

	"github.com/funthere/starset/userservice/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u userUsecase) Register(ctx context.Context, user *domain.User) error {
	return u.userRepo.Register(ctx, user)
}

func (u userUsecase) Login(user *domain.User) error {
	return u.userRepo.Login(user)
}

func (u userUsecase) GetUserById(id uint32) (domain.User, error) {
	return u.userRepo.GetUserById(id)
}

func (u userUsecase) FetchUsersByIds(ctx context.Context, ids []uint32) ([]domain.User, error) {
	return u.userRepo.FetchUsersByIds(ctx, ids)
}
