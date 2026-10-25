package producer

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type producer struct {
	producer sarama.SyncProducer
}

func New(brokers []string) (*producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("create sync producer: %w", err)
	}

	return &producer{producer: p}, nil
}

func (p *producer) Produce(topic, msg, aggregateKey string) (err error) {
	producerMsg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(aggregateKey),
		Value:     sarama.StringEncoder(msg),
		Timestamp: time.Now(),
	}

	_, _, err = p.producer.SendMessage(producerMsg)
	if err != nil {

		return fmt.Errorf("send message to topic %q: %w", topic, err)
	}

	return nil
}

func (p *producer) Close() error {

	return p.producer.Close()
}
