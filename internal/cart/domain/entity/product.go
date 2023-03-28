package entity

import (
	"time"
)

type Product struct {
	ID             string `gorm:"type:uuid;primaryKey;"`
	StoreProfileId string
	StoreProfile   StoreProfile
	Title          string
	Description    string
	UpdatedAt      time.Time
}

func (p *Product) Update(description string) {
	p.Description = description
}

func (Product) TableName() string {
	return "products"
}
