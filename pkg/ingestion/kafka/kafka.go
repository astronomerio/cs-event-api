package kafka

import (
	"fmt"
	"log"
	"strings"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaIngestionHandler struct {
	producer *kafka.Producer
	topic    string
}

var appConfig = config.Get()

func NewKafkaIngestionHandler() *KafkaIngestionHandler {
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
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

func (h *KafkaIngestionHandler) Start() error {
	h.startEventListener()
	return nil
}

func (h *KafkaIngestionHandler) Shutdown() error {
	h.producer.Flush(10 * 1000)
	h.producer.Close()
	return nil
}

func (h *KafkaIngestionHandler) startEventListener() {
	go func() {
		for e := range h.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					if appConfig.DebugMode {
						fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
							*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
					}
				}
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()
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
