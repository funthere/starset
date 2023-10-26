package usecase

import (
	"context"

	"github.com/funthere/starset/userservice/domain"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u UserUsecase) Register(ctx context.Context, user domain.User) error {
	return u.userRepo.Register(ctx, user)
}

func (u UserUsecase) Login(user *domain.User) error {
	return u.userRepo.Login(user)
}

func (u UserUsecase) GetUserById(id uint32) (domain.User, error) {
	return u.userRepo.GetUserById(id)
}

func (u *UserUsecase) FetchUsersByIds(ctx context.Context, ids []uint32) ([]domain.User, error) {
	return u.userRepo.FetchUsersByIds(ctx, ids)
}
