package repository

import "github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"

type ProductRepository interface {
	Add(product entity.Product) error
	GetByIdAndProfileId(productId string, profileId string) (*entity.Product, error)
}
