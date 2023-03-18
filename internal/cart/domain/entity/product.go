package entity

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	ID          string `gorm:"type:uuid;primaryKey;"`
	Type        string
	Title       string
	Description string
	UpdatedAt   time.Time
}

func (p *Product) Update(description string) {
	p.Description = description
}
func (Product) TableName() string {
	return "products"
}
