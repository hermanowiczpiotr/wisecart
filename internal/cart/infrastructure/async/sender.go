package async

import "github.com/nats-io/nats.go"

type NatsMessageSender struct {
	NatsConn *nats.Conn
}

func (s NatsMessageSender) Send(subject string, data []byte) error {
	return s.NatsConn.Publish(subject, data)
}
