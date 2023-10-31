package repository

import (
	"context"
	"errors"

	"github.com/funthere/starset/userservice/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Register(ctx context.Context, user *domain.User) error {
	if err := u.db.Where("email", user.Email).First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Email duplicated!")
	}

	err := u.db.Debug().Create(&user).Error

	return err
}

func (u *userRepository) Login(user *domain.User) error {
	return u.db.Where("email", user.Email).First(&user).Error
}

func (u *userRepository) GetUserById(id uint32) (domain.User, error) {
	user := domain.User{}
	if err := u.db.First(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *userRepository) FetchUsersByIds(ctx context.Context, ids []uint32) ([]domain.User, error) {
	users := []domain.User{}
	if err := u.db.Find(&users, ids).Error; err != nil {
		return []domain.User{}, err
	}
	return users, nil
}
