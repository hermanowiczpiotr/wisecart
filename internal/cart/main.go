package main

import (
	"fmt"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/commands"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/repository"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/async"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/genproto"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/persistence"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/shopify"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/ui"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/ui/cli"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"runtime"
)

func main() {
	repositories := persistence.NewRepositories()
	err := repositories.AutoMigrate()

	nc, _ := nats.Connect("nats://nats:4222")
	messageSender := async.NatsMessageSender{NatsConn: nc}

	runSubscriber(nc, repositories.ProductRepository, repositories.StoreProfileRepository)

	cliSrv := cli.ShopifyCli{
		Client:                 shopify.Client{},
		StoreProfileRepository: repositories.StoreProfileRepository,
	}

	cliSrv.Run()

	if err != nil {

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

func runSubscriber(
	nc *nats.Conn,
	productRepo repository.ProductRepository,
	storeProfileRepo repository.ProfileRepository,
) {
	log.Printf("mesydz")
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	sub, err := nc.QueueSubscribe("my.subject", "my-queue-group", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Fatal(err)
	}
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	//
	//clients := []infrastructure.ClientInterface{
	//	shopify.NewShopifyClient(),
	//}
	//clientService := service.ProductService{clients}

	//productsSubscriber := async.NewProductsSubscriber(sub, commands.SynchronizeProductsCommandHandler{
	//	productRepo, clientService, storeProfileRepo,
	//})

	//productsSubscriber.Run()

	if err != nil {
		log.Fatal(err)
	}
	//wg.Wait()
	runtime.Goexit()
}
