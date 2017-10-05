package ingestion

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion/kafka"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion/kinesis"
	"github.com/sirupsen/logrus"
)

type IngestionHandler interface {
	ProcessMessage(string, string)
	Start() error
	Shutdown() error
}

func NewHandler(kind string, log *logrus.Logger) IngestionHandler {
	logger := log.WithFields(logrus.Fields{"package": "api", "function": "NewHandler"})

	handlers := map[string]func() IngestionHandler{
		"kinesis": func() IngestionHandler {
			return kinesis.NewIngestionHandler(log)
		},
		"mock-kinesis": func() IngestionHandler {
			return kinesis.NewMockIngestionHandler()
		},
		"localstack": func() IngestionHandler {
			return kinesis.NewMockLocalStackIngestionHandler(log)
		},
		"kafka": func() IngestionHandler {
			return kafka.NewIngestionHandler(log)
		},
	}

	f, ok := handlers[kind]
	if !ok {
		logger.Fatal("invalid ingestion handler")
	}

	return f()
}
