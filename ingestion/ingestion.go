package ingestion

import (
	"github.com/astronomerio/event-api/ingestion/kafka"
	"github.com/astronomerio/event-api/ingestion/stdout"
	v1types "github.com/astronomerio/event-api/types/v1"
)

// MessageWriter is an abstract handler that should pipe events to their next destination
type MessageWriter interface {
	Write(v1types.Message) error
	Close()
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
