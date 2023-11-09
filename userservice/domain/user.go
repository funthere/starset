package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/funthere/starset/userservice/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"NOT NULL;type:varchar(255);" validate:"required"`
	Email     string    `json:"email" gorm:"NOT NULL;unique;type:varchar(255);" validate:"required,email"`
	Password  string    `json:"password,omitempty" gorm:"NOT NULL;type:varchar(255);" validate:"required"`
	Address   string    `json:"address" gorm:"NOT NULL;type:varchar(255);" validate:"required"`
	CreatedAt time.Time `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().ID()
	u.Password = helpers.HashPassword(u.Password)
	u.CreatedAt = time.Now()

	return nil
}
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID      uint32 `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Address string `json:"address"`
	}{
		ID:      u.ID,
		Email:   u.Email,
		Name:    u.Name,
		Address: u.Address,
	})
}

type UserUsecase interface {
	Register(ctx context.Context, user *User) error
	Login(user *User) error
	GetUserById(id uint32) (User, error)
	FetchUsersByIds(ctx context.Context, ids []uint32) ([]User, error)
}

type UserRepository interface {
	Register(ctx context.Context, user *User) error
	Login(user *User) error
	GetUserById(id uint32) (User, error)
	FetchUsersByIds(ctx context.Context, ids []uint32) ([]User, error)
}
