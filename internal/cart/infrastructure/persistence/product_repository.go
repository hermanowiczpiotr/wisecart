package persistence

import (
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Add(product entity.Product) error {
	err := r.db.Create(product).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) GetByIdAndProfileId(productId string, profileId string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.Where("product_id = ? and profile_id = ?", productId, profileId).Take(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}
