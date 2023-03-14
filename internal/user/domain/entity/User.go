package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string `gorm:"type:uuid;primaryKey;"`
	Email    string `gorm:"unique;index"`
	Password string
}

func (User) TableName() string {
	return "users"
}
