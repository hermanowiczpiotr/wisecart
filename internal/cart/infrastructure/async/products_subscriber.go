package async

import (
	"encoding/json"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/async"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/commands"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type ProductsSubscriber struct {
	conn    *nats.Conn
	command commands.SynchronizeProductsCommandHandler
}

func NewProductsSubscriber(conn *nats.Conn, commandHandler commands.SynchronizeProductsCommandHandler) *ProductsSubscriber {
	return &ProductsSubscriber{
		conn:    conn,
		command: commandHandler,
	}
}

func (h *ProductsSubscriber) Run() {
	workersCount := 3
	for i := 0; i < workersCount; i++ {
		_, err := h.conn.QueueSubscribe("sync_products", "my-queue-group", func(msg *nats.Msg) {
			log.Info("Received message: %s", string(msg.Data))
			if msg != nil {
				go func() {
					var storeProfileMessage async.StoreProfileAsyncMessage
					if err := json.Unmarshal(msg.Data, &storeProfileMessage); err != nil {
						log.Error("Failed to parse message: %s", err)
						return
					}

					h.command.Handle(commands.SynchronizeProductsCommand{
						StoreProfileId: storeProfileMessage.StoreProfileId,
					})

					msg.Ack()
				}()
			}
		})

		if err = h.conn.LastError(); err != nil {
			log.Fatal(err)
		}
	}
}
