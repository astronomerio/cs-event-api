package kafka

import (
	"log"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"

	"github.com/Shopify/sarama"
)

type KafkaIngestionHandler struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaIngestionHandler() *KafkaIngestionHandler {
	appConfig := config.Get()
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(appConfig.KafkaBrokers, cfg)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	h := KafkaIngestionHandler{
		producer: producer,
		topic:    appConfig.KafkaTopic,
	}
	return &h
}

func (h *KafkaIngestionHandler) ProcessMessage(r, partition string) {
	_, _, err := h.producer.SendMessage(&sarama.ProducerMessage{
		Topic: h.topic,
		Key:   sarama.StringEncoder(partition),
		Value: sarama.StringEncoder(partition),
	})
	if err != nil {
		log.Println(err)
	}
}
