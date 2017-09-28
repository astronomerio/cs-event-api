package kafka

import (
	"log"
	"strings"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaIngestionHandler struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaIngestionHandler() *KafkaIngestionHandler {
	appConfig := config.Get()
	cfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(appConfig.KafkaBrokers, ","),
	}
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create kafka client: %s\n", err)
	}
	h := KafkaIngestionHandler{
		producer: producer,
		topic:    appConfig.KafkaTopic,
	}
	return &h
}

func (h *KafkaIngestionHandler) ProcessMessage(message, partition string) {
	kmessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &h.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(partition),
		Value: []byte(message),
	}
	h.producer.ProduceChannel() <- kmessage
}
