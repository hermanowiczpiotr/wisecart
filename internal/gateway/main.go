package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hermanowiczpiotr/ola/internal/gateway/infrastructure/server"
	"github.com/hermanowiczpiotr/ola/internal/gateway/infrastructure/server/grcp"
	"github.com/hermanowiczpiotr/ola/internal/gateway/ui"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("starting server")
	cc, err := grpc.Dial(
		os.Getenv("USER_SERVICE_GRPC_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	server.StartHttpServer(
		func(router chi.Router) http.Handler {
			return server.HandlerFromMux(
				ui.NewHandler(grcp.NewUserClient(cc)),
				router,
			)
		},
	)

	log.Print("server started")
}
