package query

import (
	"github.com/hermanowiczpiotr/wisecart/internal/user/domain/entity"
	"github.com/hermanowiczpiotr/wisecart/internal/user/domain/repository"
)

type GetUserByIdQuery struct {
	Id string
}

type GetUserByIdQueryHandler struct {
	UserRepo repository.UserRepository
}

func NewGetUserByIdQuery(userRepo repository.UserRepository) GetUserByIdQueryHandler {
	return GetUserByIdQueryHandler{
		UserRepo: userRepo,
	}
}

func (userQuery GetUserByIdQueryHandler) Handle(q GetUserByIdQuery) (*entity.User, error) {
	user, err := userQuery.UserRepo.GetUserById(q.Id)

	if err != nil {
		return nil, err
	}

	return user, nil
}
