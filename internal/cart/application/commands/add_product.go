package commands

import (
	"github.com/google/uuid"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/repository"
	"time"
)

type AddProductCommand struct {
	Title       string
	Description string
}

type AddProductCommandHandler struct {
	productRepository repository.ProductRepository
}

func NewAddProductCommandHandler(productRepository repository.ProductRepository) AddProductCommandHandler {
	return AddProductCommandHandler{
		productRepository: productRepository,
	}
}

func (h AddProductCommandHandler) handle(command AddProductCommand) error {

	return h.productRepository.Add(entity.Product{
		ID:          uuid.New().String(),
		Title:       command.Title,
		Description: command.Description,
		UpdatedAt:   time.Now(),
	})
}
