package repository

import (
	"github.com/hermanowiczpiotr/wisecart/internal/user/domain/entity"
)

type UserRepository interface {
	Add(user entity.User) error
	GetUserById(id string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
}
