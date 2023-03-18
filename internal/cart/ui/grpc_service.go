package ui

import "github.com/hermanowiczpiotr/wisecart/internal/cart/application"

type GRPCService struct {
	App application.CartApp
}

func NewGRPCService(app application.CartApp) GRPCService {
	return GRPCService{
		App: app,
	}
}
