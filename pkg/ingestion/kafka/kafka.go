package kafka

import (
	"strings"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
)

type Handler struct {
	producer *kafka.Producer
	topic    string
}

var appConfig = config.Get()

func NewHandler() *Handler {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "NewHandler"})
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	cfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(appConfig.KafkaBrokers, ","),
	}
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		logger.Fatalf("failed to create kafka client: %s\n", err)
	}
	h := Handler {
		producer: producer,
		topic:    appConfig.KafkaTopic,
	}
	return &h
}

func (h *Handler) Start() error {
	h.startEventListener()
	return nil
}

func (h *Handler) Shutdown() error {
	h.producer.Flush(10 * 1000)
	h.producer.Close()
	return nil
}

func (h *Handler) startEventListener() {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "startEventListener"})
	go func() {
		for e := range h.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					logger.Infof("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					logger.Debugf("delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			default:
				logger.Infof("ignored event: %s\n", ev)
			}
		}
	}()
}

func (h *Handler) ProcessMessage(message, partition string) {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &h.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(partition),
		Value: []byte(message),
	}
	h.producer.ProduceChannel() <- msg
}
