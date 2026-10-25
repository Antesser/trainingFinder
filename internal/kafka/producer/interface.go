package producer

type MessageSender interface {
	SendMessage(topic, msg, aggregateKey string) (int32, int64, error)
	Close() error
}
