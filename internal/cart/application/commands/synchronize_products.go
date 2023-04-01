package commands

import (
	"github.com/google/uuid"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/repository"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/service"
	log "github.com/sirupsen/logrus"
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
		log.Error(err)
	}

	productsList, err := h.ProductService.GetProductsByStoreProfile(storeProfile)

	log.Error(err)

	for _, productDto := range productsList.Products {
		err = h.ProductRepository.Add(entity.Product{
			ID:             uuid.New().String(),
			StoreProfileId: storeProfile.ID,
			Title:          productDto.Title,
			Description:    productDto.Description,
			UpdatedAt:      time.Now(),
		})

		if err != nil {
			log.Error(err)
		}
	}
}
