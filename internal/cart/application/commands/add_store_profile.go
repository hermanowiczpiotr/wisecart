package commands

import (
	"github.com/google/uuid"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/repository"
)

type AddStoreProfileCommand struct {
	UserId            string
	Name              string
	Type              string
	AuthorizationData []byte
}

type AddStoreProfileCommandHandler struct {
	storeProfileRepository repository.ProfileRepository
}

func NewAddStoreProfileCommandHandler(storeProfileRepository repository.ProfileRepository) AddStoreProfileCommandHandler {
	return AddStoreProfileCommandHandler{
		storeProfileRepository: storeProfileRepository,
	}
}

func (h AddStoreProfileCommandHandler) Handle(command AddStoreProfileCommand) error {
	return h.storeProfileRepository.Add(entity.StoreProfile{
		ID:                uuid.New().String(),
		UserId:            command.UserId,
		Name:              command.Name,
		Type:              command.Type,
		AuthorizationData: command.AuthorizationData,
	})
}
