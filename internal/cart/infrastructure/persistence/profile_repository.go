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

func (r *ProfileRepository) GetByUserId(userId string) (*entity.StoreProfile, error) {
	var storeProfile entity.StoreProfile
	err := r.db.Where("user_id = ?", userId).Take(&storeProfile).Error
	if err != nil {
		return nil, err
	}

	return &storeProfile, nil
}
