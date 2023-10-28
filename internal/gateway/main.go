package main

import (
	"github.com/hermanowiczpiotr/wisecart/internal/gateway/infrastructure/logs"
	"github.com/hermanowiczpiotr/wisecart/internal/gateway/infrastructure/server"
	"github.com/hermanowiczpiotr/wisecart/internal/gateway/infrastructure/server/genproto"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func main() {
	logs.Init()
	log.Info("starting server")

	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	userClientConnection, err := grpc.Dial(
		os.Getenv("USER_SERVICE_GRPC_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
	)

	cartClientConnection, err := grpc.Dial(
		os.Getenv("CART_SERVICE_GRPC_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
	)

	if err != nil {
		log.Error("Could not connect:", err)
	}

	grcpUserClient := genproto.NewUserClient(userClientConnection)

	router := server.NewRouter(grcpUserClient, genproto.NewCartClient(cartClientConnection))

	err = router.Run(":" + os.Getenv("PORT"))

	log.Info("server started", err)
}
