package persistence

import (
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var (
	DSN = os.Getenv("cart_dsn")
)

type Repositories struct {
	db                     *gorm.DB
	ProductRepository      *ProductRepository
	StoreProfileRepository *ProfileRepository
}

func NewRepositories() *Repositories {
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &Repositories{
		db:                     db,
		ProductRepository:      NewProductRepository(db),
		StoreProfileRepository: NewProfileRepository(db),
	}
}

func (r *Repositories) AutoMigrate() error {
	return r.db.AutoMigrate(
		&entity.StoreProfile{},
		&entity.Product{},
		&entity.ProductItem{},
	)
}
