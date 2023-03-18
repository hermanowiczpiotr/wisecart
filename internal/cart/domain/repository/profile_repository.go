package repository

import "github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"

type ProfileRepository interface {
	GetByUserId(userId string) (*entity.StoreProfile, error)
}
