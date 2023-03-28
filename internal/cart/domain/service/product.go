package service

import (
	"errors"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/dto"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure"
)

type ProductService struct {
	Clients []infrastructure.ClientInterface
}

func (ps ProductService) GetProductsByStoreProfile(profile *entity.StoreProfile) (dto.ProductDtoList, error) {
	for _, client := range ps.Clients {
		if client.Support(profile) {
			return client.FetchProducts(profile)
		}
	}

	return dto.ProductDtoList{}, errors.New("not found client for provided profile")
}
