package commands

import (
	"github.com/google/uuid"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/repository"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/service"
	"log"
	"time"
)

type SynchronizeProductsCommand struct {
	StoreProfileId string
}

type SynchronizeProductsCommandHandler struct {
	ProductRepository      repository.ProductRepository
	ProductService         service.ProductService
	StoreProfileRepository repository.ProfileRepository
}

func (h SynchronizeProductsCommandHandler) Handle(command SynchronizeProductsCommand) {
	storeProfile, err := h.StoreProfileRepository.GetById(command.StoreProfileId)
	if err != nil {
		log.Printf(err.Error())
	}

	productsList, err := h.ProductService.GetProductsByStoreProfile(storeProfile)

	for _, productDto := range productsList.Products {
		h.ProductRepository.Add(entity.Product{
			ID:             uuid.New().String(),
			StoreProfileId: storeProfile.ID,
			Title:          productDto.Title,
			Description:    productDto.Description,
			UpdatedAt:      time.Now(),
		})
	}
}
