package main

import (
	"fmt"
	"github.com/hermanowiczpiotr/wisecart/internal/gateway/infrastructure/server"
	"github.com/hermanowiczpiotr/wisecart/internal/gateway/infrastructure/server/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func main() {
	log.Print("starting server")
	userClientConnection, err := grpc.Dial(
		os.Getenv("USER_SERVICE_GRPC_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	cartClientConnection, err := grpc.Dial(
		os.Getenv("CART_SERVICE_GRPC_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	grcpUserClient := genproto.NewUserClient(userClientConnection)

	router := server.NewRouter(grcpUserClient, genproto.NewCartClient(cartClientConnection))

	router.Run(":" + os.Getenv("PORT"))
	log.Print("server started")
}
