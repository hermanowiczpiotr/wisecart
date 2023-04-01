package main

import (
	"fmt"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/commands"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/service"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/async"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/genproto"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/logs"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/persistence"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/shopify"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/ui"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/ui/cli"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {

	logs.Init()
	repositories := persistence.NewRepositories()
	err := repositories.AutoMigrate()

	if err != nil {
		log.Fatal(err)
	}

	nc, _ := nats.Connect("nats://nats:4222")
	messageSender := async.NatsMessageSender{NatsConn: nc}

	productsSub := async.NewProductsSubscriber(nc, commands.SynchronizeProductsCommandHandler{
		ProductRepository:      repositories.ProductRepository,
		ProductService:         service.ProductService{},
		StoreProfileRepository: repositories.StoreProfileRepository,
	})

	productsSub.Run()

	cliSrv := cli.ShopifyCli{
		Client:                 shopify.Client{},
		StoreProfileRepository: repositories.StoreProfileRepository,
	}

	cliSrv.Run()

	if err != nil {
		log.Fatal(err)

	}

	app := application.CartApp{
		Commands: application.Commands{
			AddProductCommandHandler:      commands.NewAddProductCommandHandler(repositories.ProductRepository),
			UpdateProductCommandHandler:   commands.NewUpdateProductCommandHandler(repositories.ProductRepository),
			AddStoreProfileCommandHandler: commands.NewAddStoreProfileCommandHandler(repositories.StoreProfileRepository),
		},
	}

	runGrpcServer(ui.NewGRPCService(app, messageSender))
}

func runGrpcServer(grpcService ui.GRPCService) {
	addr := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	grpcServer := grpc.NewServer()
	genproto.RegisterCartServer(grpcServer, grpcService)

	err = grpcServer.Serve(listen)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}
}
