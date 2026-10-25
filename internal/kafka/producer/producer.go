package producer

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func New(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("create sync producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(topic, msg, aggregateKey string) (partition int32, offset int64, err error) {
	producerMsg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(aggregateKey),
		Value:     sarama.StringEncoder(msg),
		Timestamp: time.Now(),
	}

	partition, offset, err = p.producer.SendMessage(producerMsg)
	if err != nil {

		return 0, 0, fmt.Errorf("send message to topic %q: %w", topic, err)
	}

	return partition, offset, nil
}

func (p *Producer) Close() error {

	return p.producer.Close()
}
