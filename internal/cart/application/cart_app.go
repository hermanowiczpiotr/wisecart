package application

import "github.com/hermanowiczpiotr/wisecart/internal/cart/application/commands"

type CartApp struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddProductCommandHandler          commands.AddProductCommandHandler
	UpdateProductCommandHandler       commands.UpdateProductCommandHandler
	AddStoreProfileCommandHandler     commands.AddStoreProfileCommandHandler
	SynchronizeProductsCommandHandler commands.SynchronizeProductsCommandHandler
}

type Queries struct {
}
