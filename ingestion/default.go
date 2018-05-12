package ingestion

import (
	"github.com/astronomerio/event-api/logging"
	v1types "github.com/astronomerio/event-api/types/v1"
	"github.com/sirupsen/logrus"
)

// Writer simply pipes input to stdout
type Writer struct{}

// NewDefaultWriter returns a new Writer
func NewDefaultWriter() *Writer {
	return &Writer{}
}

// Write prints the message to stdout
func (h *Writer) Write(ev v1types.Message) error {
	log := logging.GetLogger(logrus.Fields{"package": "stdout"})
	log.Infof("%s ==> %s\n", ev.String(), ev.GetMessageID())
	return nil
}

// Close the handler
func (h *Writer) Close() {
	log := logging.GetLogger(logrus.Fields{"package": "stdout"})
	log.Info("Producer is closed")
}
