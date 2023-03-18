package infrastructure

import (
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/dto"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
)

type ClientInterface interface {
	FetchProducts(storeProfile *entity.StoreProfile) (dto.ProductDtoList, error)
	Support(profile *entity.StoreProfile) bool
}
