package ingestion

import (
	"github.com/arizz96/event-api/ingestion/kafka"
	v1types "github.com/arizz96/event-api/types/v1"
)

// MessageWriter is an abstract handler that should pipe events to their next destination
type MessageWriter interface {
	Write(v1types.Message) error
	Close()
}

// NewMessageWriter reads application configuration and returns a new MessageWriter
func NewMessageWriter(kind string) MessageWriter {
	if kind == "kafka" {
		return kafka.NewWriter()
	}
	return NewDefaultWriter()
}
