package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	brokers := []string{"localhost:9094"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Не удалось создать продюсер: %v", err)
	}
	defer producer.Close()

	topic := "trainings"
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(fmt.Sprintf("Дароуки, сейчас %s", time.Now().Format(time.RFC3339))),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Ошибка отправки: %v", err)
	}

	fmt.Printf("Сообщение отправлено в топик '%s', партиция %d, смещение %d\n", topic, partition, offset)
}
