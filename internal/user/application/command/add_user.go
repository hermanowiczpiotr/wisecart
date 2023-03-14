package command

import (
	"github.com/google/uuid"
	"github.com/hermanowiczpiotr/ola/internal/user/domain/entity"
	"github.com/hermanowiczpiotr/ola/internal/user/domain/repository"
)

type AddUserCommand struct {
	Email    string
	Password string
}

type AddUserCommandHandler struct {
	userRepo repository.UserRepository
}

func NewAddUserCommand(userRepo repository.UserRepository) AddUserCommandHandler {
	return AddUserCommandHandler{userRepo}
}

func (commandHandler AddUserCommandHandler) Handle(command AddUserCommand) error {

	err := commandHandler.userRepo.Add(entity.User{
		ID:       uuid.New().String(),
		Email:    command.Email,
		Password: command.Password,
	})

	if err != nil {
		return err
	}

	return nil
}
