package entity

import "gorm.io/gorm"

type ProductItem struct {
	gorm.Model
	ID        string `gorm:"type:uuid;primaryKey;"`
	ProductID string
	Product   Product
	Title     string
	Sku       string
	Price     float32
}

func (ProductItem) TableName() string {
	return "products_items"
}
