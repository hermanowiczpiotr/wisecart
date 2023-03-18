package service

import (
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure"
)

type ProductService struct {
	clients []infrastructure.ClientInterface
}

func (ps *ProductService) getProductsByStoreProfile(profile *entity.StoreProfile) {
	for _, client := range ps.clients {
		if client.Support(profile) {
			client.FetchProducts(profile)
		}
	}
}
