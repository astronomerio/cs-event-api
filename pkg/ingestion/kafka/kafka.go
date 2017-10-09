package kafka

import (
	"strings"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
)

type KafkaHandler struct {
	producer *kafka.Producer
	topic    string
}

var appConfig = config.Get()
var isRunning = false

func NewHandler() *KafkaHandler {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "NewHandler"})
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	cfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(appConfig.KafkaBrokers, ","),
	}
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		logger.Fatalf("failed to create kafka client: %s\n", err)
	}
	h := KafkaHandler {
		producer: producer,
		topic:    appConfig.KafkaTopic,
	}
	return &h
}

func (h *KafkaHandler) Start() error {
	h.startEventListener()
	return nil
}

func (h *KafkaHandler) Shutdown() error {
	h.producer.Flush(10 * 1000)
	h.producer.Close()
	return nil
}

func (h *KafkaHandler) startEventListener() {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "startEventListener"})
	go func() {
		isRunning = true
		defer func(){
			isRunning = false
		}()
		for e := range h.producer.Events(){
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					logger.Errorf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					logger.Debugf("delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			default:
				logger.Errorf("non kafka message found in event stream: %s\n", ev)
			}
		}
	}()
}

func (h *KafkaHandler) ProcessMessage(message, partition string) {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "ProcessMessage"})
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &h.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(partition),
		Value: []byte(message),
	}

	if isRunning != true {
		logger.Error("event listener isn't active")
	} else {
		err := h.producer.Produce(msg, h.producer.Events())
		if err != nil {
			logger.Errorf("unable to produce %f", err.Error())
		}
	}
}
