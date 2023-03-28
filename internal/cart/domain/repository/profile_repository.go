package repository

import "github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"

type ProfileRepository interface {
	Add(storeProfile entity.StoreProfile) error
	GetByUserId(userId string) (*entity.StoreProfile, error)
	GetById(id string) (*entity.StoreProfile, error)
}
