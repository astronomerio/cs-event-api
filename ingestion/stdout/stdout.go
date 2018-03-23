package stdout

import "fmt"

// Writer simply pipes input to stdout
type Writer struct{}

// NewWriter returns a new Writer
func NewWriter() *Writer {
	return &Writer{}
}

// Start starts the handler
func (h *Writer) Start() error {
	return nil
}

// Shutdown the handler
func (h *Writer) Shutdown() error {
	return nil
}

// ProcessMessage prints the message to stdout
func (h *Writer) ProcessMessage(message, partition string) {
	fmt.Printf("%s ==> %s\n", message, partition)
}
