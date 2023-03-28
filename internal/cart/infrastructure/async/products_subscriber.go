package async

import (
	"encoding/json"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/async"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/commands"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type ProductsSubscriber struct {
	subscription *nats.Subscription
	command      commands.SynchronizeProductsCommandHandler
}

func NewProductsSubscriber(subscription *nats.Subscription, commandHandler commands.SynchronizeProductsCommandHandler) *ProductsSubscriber {
	return &ProductsSubscriber{
		subscription: subscription,
		command:      commandHandler,
	}
}

func (h *ProductsSubscriber) Run() {

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		// Wait for messages
		msg, err := h.subscription.NextMsg(time.Second * 5)
		log.Print(msg)
		if err != nil {
			// Handle error
		}

		if msg != nil {
			var storeProfileMessage async.StoreProfileAsyncMessage
			if err := json.Unmarshal(msg.Data, &storeProfileMessage); err != nil {
				log.Printf("Failed to parse message: %s", err)
				return
			}

			h.command.Handle(commands.SynchronizeProductsCommand{
				StoreProfileId: storeProfileMessage.StoreProfileId,
			})

			err = msg.Ack()
			if err != nil {
				// Handle error
			}
		}

	}
}
