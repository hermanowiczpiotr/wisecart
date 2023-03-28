package async

type MessageSender interface {
	Send(subject string, data []byte) error
}
