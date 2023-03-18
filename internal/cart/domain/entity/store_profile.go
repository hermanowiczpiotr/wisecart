package entity

import (
	"gorm.io/gorm"
)

type StoreProfile struct {
	gorm.Model
	ID                string `gorm:"type:uuid;primaryKey;"`
	UserId            string
	Name              string
	Type              string
	AuthorizationData []byte `gorm:"type:json"`
}

func (StoreProfile) TableName() string {
	return "store_profiles"
}
