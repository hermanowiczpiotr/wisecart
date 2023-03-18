package main

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hermanowiczpiotr/wisecart/internal/user/application"
	"github.com/hermanowiczpiotr/wisecart/internal/user/application/command"
	"github.com/hermanowiczpiotr/wisecart/internal/user/application/query"
	"github.com/hermanowiczpiotr/wisecart/internal/user/infrastructure/genproto"
	"github.com/hermanowiczpiotr/wisecart/internal/user/infrastructure/logs"
	"github.com/hermanowiczpiotr/wisecart/internal/user/infrastructure/persistence"
	"github.com/hermanowiczpiotr/wisecart/internal/user/ui"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	logs.Init()
	log.Info("Starting  server")

	services := persistence.NewRepositories()
	services.AutoMigrate()

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	app := application.UserApp{
		query.NewGetUserByIdQuery(services.User),
		query.NewGetUserByEmailQuery(services.User),
		command.NewAddUserCommand(services.User),
	}

	runGrpcServer(ui.NewGRPCService(app, tokenAuth))

	log.Info("Server started")
}

func runGrpcServer(grpcService ui.GRPCService) {
	addr := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	grpcServer := grpc.NewServer()
	genproto.RegisterUserServer(grpcServer, grpcService)

	err = grpcServer.Serve(listen)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}
}
