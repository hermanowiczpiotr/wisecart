package persistence

import (
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) Add(storeProfile entity.StoreProfile) error {
	err := r.db.Create(&storeProfile).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) GetByUserId(userId string) (*entity.StoreProfile, error) {
	var storeProfile entity.StoreProfile
	err := r.db.Where("user_id = ?", userId).Take(&storeProfile).Error
	if err != nil {
		return nil, err
	}

	return &storeProfile, nil
}

func (r *ProfileRepository) GetById(id string) (*entity.StoreProfile, error) {
	var storeProfile entity.StoreProfile
	err := r.db.Where("id = ?", id).Take(&storeProfile).Error
	if err != nil {
		return nil, err
	}

	return &storeProfile, nil
}
