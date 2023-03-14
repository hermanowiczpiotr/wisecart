package query

import (
	"github.com/hermanowiczpiotr/ola/internal/user/domain/entity"
	"github.com/hermanowiczpiotr/ola/internal/user/domain/repository"
)

type GetUserByEmailQuery struct {
	Email string
}

type GetUserByEmailQueryHandler struct {
	UserRepo repository.UserRepository
}

func NewGetUserByEmailQuery(userRepo repository.UserRepository) GetUserByEmailQueryHandler {
	return GetUserByEmailQueryHandler{
		UserRepo: userRepo,
	}
}

func (userQuery GetUserByEmailQueryHandler) Handle(q GetUserByEmailQuery) (*entity.User, error) {
	user, err := userQuery.UserRepo.GetUserByEmail(q.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}
