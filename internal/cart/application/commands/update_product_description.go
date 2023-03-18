package commands

import (
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/repository"
)

type UpdateProductCommand struct {
	ProductId   string
	StoreId     string
	Description string
}

type UpdateProductCommandHandler struct {
	productRepository repository.ProductRepository
}

func NewUpdateProductCommandHandler(productRepository repository.ProductRepository) UpdateProductCommandHandler {
	return UpdateProductCommandHandler{
		productRepository: productRepository,
	}
}

func (h UpdateProductCommandHandler) handle(command UpdateProductCommand) error {
	product, err := h.productRepository.GetByIdAndProfileId(command.ProductId, command.StoreId)

	if err != nil {
		return err
	}

	product.Update(command.Description)

	return nil
}
