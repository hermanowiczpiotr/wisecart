package persistence

import (
	"github.com/hermanowiczpiotr/wisecart/internal/user/domain/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Add(user entity.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserById(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).Take(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).Take(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
