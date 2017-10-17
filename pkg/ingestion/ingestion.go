package ingestion

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion/kafka"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion/kinesis"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	ProcessMessage(string, string)
	Start() error
	Shutdown() error
}

func NewHandler(kind string) Handler {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "api", "function": "NewHandler"})

	handlers := map[string]func() Handler{
		"kinesis": func() Handler {
			return kinesis.NewHandler()
		},
		"mock-kinesis": func() Handler {
			return kinesis.NewMockHandler()
		},
		"localstack": func() Handler {
			return kinesis.NewMockLocalStackHandler()
		},
		"kafka": func() Handler {
			return kafka.NewHandler()
		},
	}

	f, ok := handlers[kind]
	if !ok {
		logger.Fatal("invalid ingestion handler")
	}

	return f()
}
