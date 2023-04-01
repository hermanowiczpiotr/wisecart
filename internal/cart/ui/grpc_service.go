package ui

import (
	"context"
	"encoding/json"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/async"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/commands"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/genproto"
	"log"
	"net/http"
)

type GRPCService struct {
	App       application.CartApp
	Publisher async.MessageSender
}

func NewGRPCService(app application.CartApp, publisher async.MessageSender) GRPCService {
	return GRPCService{
		App:       app,
		Publisher: publisher,
	}
}

func (gs GRPCService) AddProfile(context context.Context, addStoreProfileRequest *genproto.AddStoreProfileRequest) (*genproto.AddStoreProfileResponse, error) {
	gs.App.Commands.AddStoreProfileCommandHandler.Handle(
		commands.AddStoreProfileCommand{
			UserId:            addStoreProfileRequest.UserId,
			Name:              addStoreProfileRequest.Name,
			Type:              addStoreProfileRequest.Type,
			AuthorizationData: addStoreProfileRequest.AuthorizationData,
		})

	return &genproto.AddStoreProfileResponse{
		Status: http.StatusOK,
	}, nil
}

func (gs GRPCService) SynchronizeProducts(context context.Context, synchronizeProductsRequest *genproto.SynchronizeProductsRequest) (*genproto.SynchronizeProductsResponse, error) {
	log.Printf("duuuuupa")
	data, err := json.Marshal(async.StoreProfileAsyncMessage{synchronizeProductsRequest.GetProfileId()})
	if err != nil {
		return &genproto.SynchronizeProductsResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, err
	}

	gs.Publisher.Send("sync_products", data)
	return &genproto.SynchronizeProductsResponse{
		Status: http.StatusOK,
		Error:  "",
	}, nil
}
