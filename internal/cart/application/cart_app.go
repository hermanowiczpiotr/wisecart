package application

import "github.com/hermanowiczpiotr/wisecart/internal/cart/application/commands"

type CartApp struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddProductCommandHandler    commands.AddProductCommandHandler
	UpdateProductCommandHandler commands.UpdateProductCommandHandler
}

type Queries struct {
}
