package ingestion

import (
	"github.com/astronomerio/event-api/ingestion/kafka"
	"github.com/astronomerio/event-api/ingestion/stdout"
)

// MessageWriter is an abstract handler that should pipe events to their next destination
type MessageWriter interface {
	Start() error
	ProcessMessage(string, string)
	Shutdown() error
}

// NewMessageWriter reads application configuration and returns a new MessageWriter
func NewMessageWriter(kind string) MessageWriter {
	switch kind {
	case "kafka":
		return kafka.NewWriter()
	default:
		return stdout.NewWriter()
	}
}
