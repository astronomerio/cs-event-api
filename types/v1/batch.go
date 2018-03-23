package v1

import (
	"encoding/json"
	"time"
)

// Batch defines the payload of a request to /batch
// analytics-go does not export this type
type Batch struct {
	MessageID    string       `json:"messageId,omitempty"`
	SentAt       time.Time    `json:"sentAt,omitempty"`
	Messages     Messages     `json:"batch,omitempty"`
	Context      *Context     `json:"context,omitempty"`
	Integrations Integrations `json:"integrations,omitempty"`
}

// Messages is a colleciton of messages.
type Messages []Message

// UnmarshalJSON turns raw bytes into messages.
func (m *Messages) UnmarshalJSON(data []byte) error {
	tmp := []message{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	var msgs = []Message{}
	for _, msg := range tmp {
		newMsg, err := NewMessage(msg.Type, msg.Raw)
		if err != nil {
			return err
		}
		msgs = append(msgs, newMsg)
	}

	*m = msgs
	return nil
}

// message is an intermediate type used to unmarshl the raw bytes into the correct type.
type message struct {
	Type string `json:"type"`
	Raw  []byte
}

// UnmarshalJSON turns raw bytes into a message.
func (m *message) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type string `json:"type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	m.Type = tmp.Type
	m.Raw = data
	return nil
}
