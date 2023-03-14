package application

import (
	"github.com/hermanowiczpiotr/ola/internal/user/application/command"
	"github.com/hermanowiczpiotr/ola/internal/user/application/query"
)

type UserApp struct {
	GetUserByIdQueryHandler    query.GetUserByIdQueryHandler
	GetUserByEmailQueryHandler query.GetUserByEmailQueryHandler
	AddUserCommandHandler      command.AddUserCommandHandler
}
