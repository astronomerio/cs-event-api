package ingestion

import (
	"log"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion/kafka"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion/kinesis"
)

type IngestionHandler interface {
	ProcessMessage(string, string)
}

func NewHandler(kind string) IngestionHandler {
	handlers := map[string]func() IngestionHandler{
		"kinesis": func() IngestionHandler {
			return kinesis.NewKinesisIngestionHandler()
		},
		"mock-kinesis": func() IngestionHandler {
			return kinesis.NewMockKinesisIngestionHandler()
		},
		"localstack": func() IngestionHandler {
			return kinesis.NewMockKinesisLocalStackIngestionHandler()
		},
		"kafka": func() IngestionHandler {
			return kafka.NewKafkaIngestionHandler()
		},
	}

	f, ok := handlers[kind]
	if !ok {
		log.Fatal("invalid ingestion handler")
	}

	return f()
}
