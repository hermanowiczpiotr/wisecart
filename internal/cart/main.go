package main

import (
	"fmt"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/persistence"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/shopify"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/ui"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/ui/cli"

	//"github.com/hermanowiczpiotr/wisecart/internal/user/infrastructure/genproto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	repositories := persistence.NewRepositories()
	err := repositories.AutoMigrate()

	cli := cli.ShopifyCli{
		Client:                 shopify.Client{},
		StoreProfileRepository: repositories.StoreProfileRepository,
	}

	cli.Run()
	if err != nil {

	}
	//
	//app := application.CartApp{
	//	Commands: application.Commands{
	//		AddProductCommandHandler: commands.NewAddProductCommandHandler(repositories.ProductRepository),
	//		UpdateProductCommandHandler: commands.NewUpdateProductCommandHandler(repositories.ProductRepository),
	//	},
	//}

	//runGrpcServer(ui.NewGRPCService(app))
}

func runGrpcServer(grpcService ui.GRPCService) {
	addr := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	grpcServer := grpc.NewServer()
	//genproto.RegisterUserServer(grpcServer, grpcService)

	err = grpcServer.Serve(listen)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}
}
