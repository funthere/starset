package domain

import (
	"context"
	"time"

	"github.com/funthere/starset/userservice/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"NOT NULL;type:varchar(255);"`
	Email     string    `json:"email" gorm:"NOT NULL;unique;type:varchar(255);"`
	Password  string    `json:"password,omitempty" gorm:"NOT NULL;type:varchar(255);"`
	Address   string    `json:"address" gorm:"NOT NULL;type:varchar(255);"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().ID()
	u.Password = helpers.HashPassword(u.Password)
	u.CreatedAt = time.Now()

	return nil
}

type UserUsecase interface {
	Register(ctx context.Context, user User) error
	Login(user *User) error
	GetUserById(id uint32) (User, error)
	FetchUsersByIds(ctx context.Context, ids []uint32) ([]User, error)
}

type UserRepository interface {
	Register(ctx context.Context, user User) error
	Login(user *User) error
	GetUserById(id uint32) (User, error)
	FetchUsersByIds(ctx context.Context, ids []uint32) ([]User, error)
}
